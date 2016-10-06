package ducklib

import (
	"github.com/carneades/carneades-4/src/engine/caes"
	"github.com/carneades/carneades-4/src/engine/terms"
)

type BoolValue struct {
	Value   bool `json:"value"`
	Assumed bool `json:"assumed"` // if true assumed, otherwise proven
}

type StmtExplanation struct {
	ConsentRequired   BoolValue `json:"consentRequired"`  // informed consent required
	Pii               BoolValue `json:"pii"`              // personally identifiable information
	Li                BoolValue `json:"li"`               // legitimate interest in the pii
	CompatiblePurpose []string `json:"compatiblePurpose"` // ids of statements with a proven compatible purpose
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

// ConsentRequired
func cr(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	// find the notConsentRequired statement for this dus
	goal, ok := terms.ReadString("notConsentRequired(" + dus.String() + ")")
	if !ok {
		return BoolValue{true, true}
	}
	for wff, stmt := range ag.Statements {
		t, ok := terms.ReadString(wff)
		var b terms.Bindings
		if !ok {
			continue
		}
		_, ok = terms.Match(goal, t, b)
		if ok {
			v := stmt.Label == caes.In // is notConsentRequired In?
			return BoolValue{!v, !v}   // consentRequired assumed
		}
	}
	return BoolValue{true, true} // consentRequired assumed
}

// Personally Identifiable Information
func pii(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	// find the notPii statement for this dus
	goal, ok := terms.ReadString("notPii(" + dus.String() + ")")
	if !ok {
		return BoolValue{true, true}
	}
	for wff, stmt := range ag.Statements {
		t, ok := terms.ReadString(wff)
		var b terms.Bindings
		if !ok {
			continue
		}
		_, ok = terms.Match(goal, t, b)
		if ok {
			v := stmt.Label == caes.In // is notPii In?
			return BoolValue{!v, !v}   // pii assumed
		}
	}
	return BoolValue{true, true} // pii assumed
}

// Legitimate Interest
func li(dus terms.Compound, ag *caes.ArgGraph) BoolValue {
	// find the li statement for this dus
	goal, ok := terms.ReadString("li(" + dus.String() + ")")
	if !ok {
		return BoolValue{true, true}
	}
	for wff, stmt := range ag.Statements {
		t, ok := terms.ReadString(wff)
		var b terms.Bindings
		if !ok {
			continue
		}
		_, ok = terms.Match(goal, t, b)
		if ok {
			v := stmt.Label == caes.In // is li In?
			return BoolValue{!v, !v}   // notLi assumed
		}
	}
	return BoolValue{false, true} // notLi assumed
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
