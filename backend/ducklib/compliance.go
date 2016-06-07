package ducklib

import (
	"io"

	"github.com/carneades/carneades-4/src/engine/caes"
)

type ArgMapFormat int

const (
	SVG ArgMapFormat = iota
	PNG
	GraphML
)

// A ComplianceChecker manages communication between the DUCK Web Server and
// the Carneades Argumentation System

type ComplianceChecker interface {
	// GetTheory: Retrieve the theory for the given ruleBaseId. If the
	// revision is not equal to the revision used to compile the theory,
	// the theory is first updated, by downloading the given revision from
	// the document database and compiling the rulebase into a theory.
	// If there are no errors, the returned error will be nil.
	GetTheory(db Database, ruleBaseId string, revision string) (caes.Theory, error)

	/* Check does the following:
		* Reads the data use document from its given io.Reader
		* Translates the data use statements in the document into Carneades assumptions (terms)
	    * Applies the theory to the assumptions, using the Carneades inference engine,
	      to construct a Carneades argument graph
		* Evaluates the argument graph to label the statements in the graph in, out or undecided.
		* Returns the evaluated argument graph
		If there are not errors, nil is returned.  Otherwise an error is returned
		describing the error.
	*/
	Check(ruleBase *caes.Theory, document *Document) (*caes.ArgGraph, error)

	isCompliant(ag *caes.ArgGraph) bool
	nonCompliantDataUseStatements(ag *caes.ArgGraph) []Statement
	Render(ag *caes.ArgGraph, format ArgMapFormat, w io.Writer) error
}
