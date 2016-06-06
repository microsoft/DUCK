package ducklib

import (
	"io"
	// "github.com/carneades/carneades-4/"
)

// A ComplianceChecker manages communication between the DUCK Web Server and
// the Carneades Argumentation System

type ComplianceChecker interface {
	/* Check does the following:
		   * Reads the rulebase from its io.Reader
		   * Reads the data use document from its given io.Reader
		   * Translates the data use statements in the document into Carneades assumptions (terms)
		   * Compiles the rulebase into a Carneades theory
	       * Applies the theory to the assumptions, using the Carneades inference engine,
	         to construct a Carneades argument graph
		   * Evaluates the argument graph to label the statements in the graph in, out or undecided.
		   * Exports the argument graph to SVG, by writing the SVG to the given io.Writer
		  If there are not errors, nil is returned.  Otherwise an error is returned
		  describing the error.
	*/
	Check(ruleBase io.Reader, dataUseDocument io.Reader, w io.Writer) error
}
