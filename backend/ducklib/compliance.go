package ducklib

import (
	// "io"

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
	GetTheory(db *Database, ruleBaseId string, revision string) (*caes.Theory, error)

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
	IsCompliant(ruleBase *caes.Theory, document *Document) (bool, error)

	/*
		CompliantDocuments does the following:
			* Translates the data use statements in the document into Carneades assumptions (terms)
			* Applies the theory to the assumptions, using the Carneades inference engine,
			  to construct a Carneades argument graph
		    * Evaluates the argument graph to label the statements in the graph in, out or undecided.
			* Returns a channel of pointers to compliant data use documents based on the
			  input document.  If the input document is compliant, a pointer to it will be returned,
			  and it will be the only document returned in the channel. If the input document is not
			  compliant, the documents returned in the channel are based on the input document, with
			  minimal changes sufficient to achieve compliance. the input document is not modified.
			  If nil is read from the channel, it is empty and contains no further references to documents.
		The error returned will be nil if and only if no errors occur this process.
	*/
	CompliantDocuments(ruleBase *caes.Theory, document *Document) (chan *Document, error)
}
