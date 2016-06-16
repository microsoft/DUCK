package ducklib

import (
	"fmt"
	"io"
	"os"
	// "path/filepath"
)

type RuleBaseDescription struct {
	Filename    string
	Id          string
	Version     string
	Title       string
	Description string
}

type ComplianceCheckerPlugin struct {
	checker     *ComplianceChecker
	RuleBaseDir string
	RuleBases   map[string]RuleBaseDescription // RuleBaseDescription.Id -> RuleBaseDescription
}

// MakeComplianceCheckerPlugin returns an error if the ruleBase dir does not
// exist or is not readable.
func MakeComplianceCheckerPlugin(ruleBaseDir string) (*ComplianceCheckerPlugin, error) {
	//	f, e := os.Open(ruleBaseDir)
	//	if e != nil {
	//		return nil, e
	//	}
	i, e := os.Stat(ruleBaseDir)
	if e != nil {
		return nil, e
	}
	if !i.IsDir() {
		return nil, fmt.Errorf("ruleBaseDir %s is not a directory", ruleBaseDir)
	}
	return &ComplianceCheckerPlugin{MakeComplianceChecker(), ruleBaseDir, make(map[string]RuleBaseDescription)}, nil
}

// ComplianceCheckerPlugin.Init: For each file in RuleBaseDir:
//    1. Parse the YAML and extract the id, version, title and description
//    2. Call checker.GetTheory function to compile each rulebase into a
//       Carneades theory and cache the theory.
//    3. Create a RuleBaseDescription and add it to the RuleBases map, indexed
//       by its Id.
//  Return an error if any rulebase cannot be compiled into a Theory
func (c ComplianceCheckerPlugin) Init() error {
	// ToDo
	return nil
}

func (c ComplianceCheckerPlugin) Shutdown() {
	// Nothing to do
}

// ruleBaseReader: returns a reader for reading the JSON source file of the
// rulebase with the given id
func ruleBaseReader(ruleBaseId string) io.Reader {
	// ToDo
	return nil
}

// IsCompliant: returns true iff the document complies with the rules in the given
// rulebase.  An error is returned if document has syntax errors and cannot be parsed.
func (c ComplianceCheckerPlugin) IsCompliant(ruleBaseId string, document *Document) (bool, error) {
	r := ruleBaseReader(ruleBaseId)
	theory, err := c.checker.GetTheory(ruleBaseId, "irrelevant", r)
	if err != nil {
		return false, err
	}
	return c.checker.IsCompliant(theory, document)
}

// CompliantDocuments: returns true iff the document complies with the rules in the given
// rulebase.  An error is returned if document has syntax errors and cannot be parsed. If
// the document is not compliant, false is returned along with a slice of compliant documents
// based on the input document. At most maxResults documents are returned. If offset is greater than
// 0, the first offset compliant documents found are skipped, allowing CompliantDocuments to be
// called repeatedly to scroll through all compliant documents incrementally.  The search
// for compliant documents is restarted each time CompliantDocuments is called, no matter
// what the offset is.
func (c ComplianceCheckerPlugin) CompliantDocuments(ruleBaseId string, document *Document, maxResults int, offset int) (bool, []*Document, error) {
	//r := ruleBaseReader(ruleBaseId)
	//theory, err := c.checker.GetTheory(ruleBaseId, "irrelevant", r)
	//if err != nil {
	//	return false, nil, err
	//}
	//cncl := MakeCanceller()
	//compliant, docs, err := c.checker.CompliantDocuments(theory, document, cncl)
	//if compliant {
	//	return true, nil, nil
	//}
	// ToDo
	// Get at most maxResults from the docs channel, skippint the offset amount and then
	// call c.Cancel() to cancel the search for other compliant documents and
	// free the resourcs of the coroutine.
	return true, nil, nil
}
