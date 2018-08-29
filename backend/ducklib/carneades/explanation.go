// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package carneades

import (
	"github.com/carneades/carneades-4/src/engine/caes"
	"github.com/carneades/carneades-4/src/engine/terms"
)

//BoolValue explains a field in the StmtExplanation in more detail.
// Value is true, when the field applies at all
//assumed is false, when the value of the Value field could be proven by carneades
type BoolValue struct {
	Value   bool `json:"value"`
	Assumed bool `json:"assumed"` // if true assumed, otherwise proven
}

//StmtExplanation represents the Explanation for one Statement.
//If the original statement has one or more and or except clauses,
//a StmtExplanation represents only one of these clauses
type StmtExplanation struct {
	ConsentRequired             BoolValue `json:"consentRequired"`             // informed consent required
	Pii                         BoolValue `json:"pii"`                         // personally identifiable information
	Li                          BoolValue `json:"li"`                          // legitimate interest in the pii
	CompatiblePurpose           []string  `json:"compatiblePurpose"`           // ids of statements with a proven compatible purpose
	IDNotRequired               BoolValue `json:"idNotRequired"`               // identification of data subject is not required
	TransferPii                 BoolValue `json:"transferPii"`                 // is there a location for a scope ans pii
	ConsentRequired2TransferPii BoolValue `json:"consentRequired2TransferPii"` // informed consent required to transfer pii
}

//Explanation contains the StmtExplanantions for each statement
// keys are statement tracking ids
type Explanation map[string]StmtExplanation

func isDataUseStatement(t terms.Term) bool {
	t2, ok := t.(terms.Compound)
	if !ok {
		return false
	}
	p, _ := terms.Predicate(t2)
	if p == "dataUseStatement" {
		return true
	}
	return false

}

// stmtId: selects the id in a DUS term and returns
// it as a string
func stmtID(t terms.Compound) string {
	return t.Args[9].String()
}

// isTrue: check whether a given predicate is true/in for a particular
// data use statement in the argument graph
func isTrue(predicate string, dus terms.Compound, ag *caes.ArgGraph, defaultValue bool) BoolValue {
	goal, ok := terms.ReadString(predicate + "(" + dus.String() + ")")
	if !ok {
		//fmt.Println("Improper predicate: " + predicate)
		return BoolValue{defaultValue, true}
	}
	for wff, stmt := range ag.Statements {
		t, ok := terms.ReadString(wff)
		var b terms.Bindings
		if !ok {
			continue
		}
		_, ok = terms.Match(goal, t, b)
		if ok {
			v := stmt.Label == caes.In
			//negation as failure, so no longer an assumption, no matter the value of v
			return BoolValue{v, false}
		}
	}
	//fmt.Println("No match found: " + predicate + "(" + dus.String() + ")")
	return BoolValue{defaultValue, true}
}

// ConsentRequired
func cr(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	return isTrue("consentRequired", dus, ag, true)
}

// Personally Identifiable Information
func pii(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	return isTrue("pii", dus, ag, true)
}

// Legitimate Interest
func li(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	return isTrue("li", dus, ag, false)
}

// Returns a slice of statements ids, for statements having a purpose
// compatible with the statement represented by wff
func cp(dus terms.Compound, ag *caes.ArgGraph) []string {
	// find the li statement for this dus
	goal, ok := terms.ReadString("compatiblePurpose(" + dus.String() + ",X)")
	if !ok {
		return []string{}
	}
	result := []string{}
	for wff, stmt := range ag.Statements {
		t, ok := terms.ReadString(wff)
		var b terms.Bindings
		if !ok {
			continue
		}
		_, ok = terms.Match(goal, t, b)
		if ok {
			if stmt.Label == caes.In {
				dus2 := t.(terms.Compound).Args[1].(terms.Compound)
				result = append(result, stmtID(dus2))
			}
		}
	}
	return result
}

// Returns a slice of ids for statements not
// requring the identification of the data subject
func idnr(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	return isTrue("idNotRequired", dus, ag, false)
}

func tpii(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	return isTrue("transferPii", dus, ag, false)
}

func cr2tpii(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	return isTrue("consentRequired2TransferPii", dus, ag, false)
}

//GetExplanation returns the Explanation struct filled with explanations for each Statement
func (c ComplianceChecker) GetExplanation(theory *caes.Theory, ag *caes.ArgGraph) (Explanation, error) {
	m := make(map[string]StmtExplanation)
	for wff := range ag.Statements {
		t, ok := terms.ReadString(wff)
		if !ok {
			continue
		}

		if isDataUseStatement(t) {
			dus := t.(terms.Compound).Args[0].(terms.Compound)
			m[stmtID(dus)] = StmtExplanation{
				ConsentRequired:             cr(dus, ag),
				Pii:                         pii(dus, ag),
				Li:                          li(dus, ag),
				CompatiblePurpose:           cp(dus, ag),
				IDNotRequired:               idnr(dus, ag),
				TransferPii:                 tpii(dus, ag),
				ConsentRequired2TransferPii: cr2tpii(dus, ag),
			}
		}
	}
	// Begin DEBUG
	//	b, err := json.Marshal(m)
	//	if err != nil {
	//		fmt.Println("error:", err)
	//	} else {
	//		// os.Stdout.Write(b)
	//		fmt.Printf("%s\n", b)
	//	}
	// End DEBUG
	return m, nil
}
