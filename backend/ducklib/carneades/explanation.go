package carneades

import (
	"github.com/carneades/carneades-4/src/engine/caes"
	"github.com/carneades/carneades-4/src/engine/terms"
)

type BoolValue struct {
	Value   bool `json:"value"`
	Assumed bool `json:"assumed"` // if true assumed, otherwise proven
}

type StmtExplanation struct {
	ConsentRequired   BoolValue `json:"consentRequired"`   // informed consent required
	Pii               BoolValue `json:"pii"`               // personally identifiable information
	Li                BoolValue `json:"li"`                // legitimate interest in the pii
	CompatiblePurpose []string  `json:"compatiblePurpose"` // ids of statements with a proven compatible purpose
	IdNotRequired     BoolValue `json:"idNotRequired"`     // identification of data subject is not required
}

type Explanation map[string]StmtExplanation // keys are statement tracking ids

func isDataUseStatement(t terms.Term) bool {
	t2, ok := t.(terms.Compound)
	if !ok {
		return false
	}
	p, _ := terms.Predicate(t2)
	if p == "dataUseStatement" {
		return true
	} else {
		return false
	}
}

// stmtId: selects the id in a DUS term and returns
// it as a string
func stmtId(t terms.Compound) string {
	return t.Args[6].String()
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
				result = append(result, stmtId(dus2))
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

func (c ComplianceChecker) GetExplanation(theory *caes.Theory, ag *caes.ArgGraph) (Explanation, error) {
	m := make(map[string]StmtExplanation)
	for wff, _ := range ag.Statements {
		t, ok := terms.ReadString(wff)
		if !ok {
			continue
		}
		if isDataUseStatement(t) {
			dus := t.(terms.Compound).Args[0].(terms.Compound)
			m[stmtId(dus)] = StmtExplanation{
				ConsentRequired:   cr(dus, ag),
				Pii:               pii(dus, ag),
				Li:                li(dus, ag),
				CompatiblePurpose: cp(dus, ag),
				IdNotRequired:     idnr(dus, ag),
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
