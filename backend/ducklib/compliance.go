package ducklib

import (
	"fmt"
	"io"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/carneades/carneades-4/src/engine/caes"
	y "github.com/carneades/carneades-4/src/engine/caes/encoding/yaml"
)

type Canceller chan struct{}

func MakeCanceller() Canceller {
	return make(chan struct{})
}

func (c Canceller) Cancel() {
	close(c)
}

type VersionedTheory struct {
	revision string
	theory   *caes.Theory
}

type ComplianceChecker struct {
	Theories map[string]VersionedTheory
}

func MakeComplianceChecker() *ComplianceChecker {
	return &ComplianceChecker{make(map[string]VersionedTheory)}
}

// GetTheory: Retrieve the theory for the given ruleBaseId. If no version of the
// rulebase has been compiled or its revision is not equal to the revision
// used to compile the theory,
// the theory is first updated, by reading the JSON source from rbSrc,
// and compiling the rulebase into a theory and updating the Theories of the
// ComplianceChecker.
// If there are no errors, the returned error will be nil.
func (c ComplianceChecker) GetTheory(ruleBaseId string, revision string, rbSrc io.Reader) (*caes.Theory, error) {
	vt, found := c.Theories[ruleBaseId]
	if !found || revision != vt.revision {

		// Compile the rulebase, update the theory cache and return the
		// theory.  Or return an error if the rulebase cannot be compiled.
		ag, err := y.Import(rbSrc)
		if err != nil {
			return nil, err
		}
		c.Theories[ruleBaseId] = VersionedTheory{revision, ag.Theory}
		return ag.Theory, nil
	}
	return vt.theory, nil

}

/*
	IsCompliant does the following:
		* Translates the data use statements in the document into Carneades assumptions (terms)
		* Applies the theory to the assumptions, using the Carneades inference engine,
		    to construct a Carneades argument graph
	    * Evaluates the argument graph to label the statements in the graph in, out or undecided.
		* Returns true if and only if the statement in the argument representing
		    the proposition that the document is compliant is in.
	The error returned will be nil if and only if no errors occur this process.
*/
func (c ComplianceChecker) IsCompliant(theory *caes.Theory, document *structs.Document) (bool, error) {
	// Construct the argument graph
	compliant := &caes.Statement{Id: "compliant",
		Metadata: make(map[string]interface{}),
		Text:     "The document is compliant.",
		Args:     []*caes.Argument{}}
	notCompliant := &caes.Statement{Id: "¬compliant",
		Metadata: make(map[string]interface{}),
		Text:     "The document is not compliant.",
		Args:     []*caes.Argument{}}
	ag := caes.ArgGraph{
		Theory:      theory,
		Assumptions: make(map[string]bool),
		Statements: map[string]*caes.Statement{
			"compliant":  compliant,
			"¬compliant": notCompliant,
		},
		Issues: map[string]*caes.Issue{
			"i1": &caes.Issue{
				Id:        "i1",
				Metadata:  make(map[string]interface{}),
				Positions: []*caes.Statement{compliant, notCompliant}},
		}}
	// add statements for the data use statements in the document
	// to the argument graph, and assume them to be true.
	for _, s := range document.Statements {
		var passive bool
		if s.Passive {
			passive = true
		} else {
			passive = false
		}
		stmtId := fmt.Sprintf("dataUseStatement(dus(%s,%s,%s,%s,%s,%s,%s,%t))",
			s.UseScopeCode,
			s.QualifierCode,
			s.DataCategoryCode,
			s.SourceScopeCode,
			s.ActionCode,
			s.ResultScopeCode,
			s.TrackingID,
			passive)

		stmt := &caes.Statement{
			Id:       stmtId,
			Metadata: make(map[string]interface{}),
			Text:     stmtId,
			Args:     []*caes.Argument{}}
		ag.Assumptions[stmtId] = true
		ag.Statements[stmtId] = stmt
	}
	// derive arguments by applying the theory of the argument graph to
	// its assumptions

	err := ag.Infer()
	if err != nil {
		return false, err
	}

	// evaluate the argument graph
	l := ag.GroundedLabelling()

	// return true iff the compliance statement is in
	return l[compliant] == caes.In, nil
}

func removeStatement(d *structs.Document, i int) (*structs.Document, error) {
	if i < 0 || i > len(d.Statements)-1 {
		return nil, fmt.Errorf("Statements index out of bounds: %v", i)
	}
	// copy the document struct
	d2 := *d
	// replace the statements with a copy, but with the selected
	// statement removed.
	d2.Statements = []structs.Statement{}
	for j, _ := range d.Statements {
		if i != j {
			d2.Statements = append(d2.Statements, d.Statements[j])
		}
	}
	return &d2, nil
}

/*
	CompliantDocuments does the following:
		* Translates the data use statements in the document into Carneades assumptions (terms)
		* Applies the theory to the assumptions, using the Carneades inference engine,
		  to construct a Carneades argument graph
	    * Evaluates the argument graph to label the statements in the graph in, out or undecided.
		* Starts a coroutine to search for compliant data use documents and returns a channel of pointers
		  to the compliant documents found. If the input document is compliant, the bool result will be true
		  and the channel returned will be closed. If the input document is not
		  compliant, the bool result will be false, and compliant alternative documents based in
		  input document will returned in the channel. The documents returend will have
		  minimal changes sufficient to achieve compliance. The input document is not modified.
		  The coroutine closes the channel when it has finished the search for compliant documents.
	An error will be returned only if was not possible to check the compliance of the input document,
	before starting the coroutine to search for compliant alternatives.
	The caller must bind c to a newly constructed Canceller, with MakeCanceller().
	If no error is returned (i.e. error is nil) the caller should call c.Cancel() when no further
	documents are needed, to cause the coroutine to be terminated.

*/
func (c ComplianceChecker) CompliantDocuments(theory *caes.Theory, doc *structs.Document, cncl Canceller) (bool, <-chan *structs.Document, error) {
	compliant, err := c.IsCompliant(theory, doc)
	if err != nil {
		return false, nil, err
	}

	if compliant {
		return true, nil, nil
	}

	// Traverse the space of documents, breadth-first, with subsets of the data use
	// statements of doc, and push each alternative documents down the
	// variants channel.

	variants := make(chan *structs.Document)

	// generate documents with subsets of the data use documents
	// and push them down the varients channel
	var subsets func(int, *structs.Document)
	subsets = func(i int, d *structs.Document) {
		if i < 0 {
			return
		}
		select {
		case <-cncl:
			return // cancelled
		default: // continue
		}
		// ToDo: check that the statement order below
		// performs breadth-first search, not depth-first
		d2, err := removeStatement(d, i)
		if err != nil {
			return
		} // indexing error should not happen
		variants <- d2
		subsets(i-1, d)
		subsets(i-1, d2)
	}
	go subsets(len(doc.Statements)-1, doc)

	// Filter out the noncompliant documents
	compliantVariants := make(chan *structs.Document)

	go func() {
		for {
			select {
			case <-cncl:
				return // cancelled
			case d3, ok := <-variants:
				if !ok { // no more variants available
					return
				}

				compliant, err := c.IsCompliant(theory, d3)
				if err != nil && compliant {
					compliantVariants <- d3
				}
			default: // do nothing, to avoid blocking
			}
		}
	}()

	return false, compliantVariants, nil
}
