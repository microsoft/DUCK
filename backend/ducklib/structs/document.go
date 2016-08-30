package structs

type Document struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Revision      string      `json:"revision"`
	Owner         string      `json:"owner"`
	Locale        string      `json:"locale"`
	AssumptionSet string      `json:"assumptionSet"`
	Statements    []Statement `json:"statements"`
}

type Statement struct {
	UseScopeCode     string `json:"useScopeCode"`
	QualifierCode    string `json:"qualifierCode"`
	DataCategoryCode string `json:"dataCategoryCode"`
	SourceScopeCode  string `json:"sourceScopeCode"`
	ActionCode       string `json:"actionCode"`
	ResultScopeCode  string `json:"resultScopeCode"`
	TrackingID       string `json:"trackingId"`
	Passive          bool   `json:"passive"`
}

func (d *Document) FromValueMap(mp map[string]interface{}) {

	if id, ok := mp["_id"]; ok {
		d.ID = id.(string)
	}
	if rev, ok := mp["_rev"]; ok {
		d.Revision = rev.(string)
	}
	if name, ok := mp["name"]; ok {
		d.Name = name.(string)
	}
	if owner, ok := mp["owner"]; ok {
		d.Owner = owner.(string)
	}
	if locale, ok := mp["locale"]; ok {
		d.Locale = locale.(string)
	}
	if assumptionSet, ok := mp["assumptionSet"]; ok {
		d.AssumptionSet = assumptionSet.(string)
	}

	d.Statements = make([]Statement, 0)
	if stmts, prs := mp["statements"].([]interface{}); prs {

		for _, stmt := range stmts {
			s := new(Statement)
			if stmt != nil {
				s.FromInterfaceMap(stmt.(map[string]interface{}))
				d.Statements = append(d.Statements, *s)
			}

		}
	}

}

func (s *Statement) FromInterfaceMap(mp map[string]interface{}) {

	s.UseScopeCode = getFieldValue(mp, "useScopeCode")
	s.QualifierCode = getFieldValue(mp, "qualifierCode")
	s.DataCategoryCode = getFieldValue(mp, "dataCategoryCode")
	s.SourceScopeCode = getFieldValue(mp, "sourceScopeCode")
	s.ActionCode = getFieldValue(mp, "actionCode")
	s.ResultScopeCode = getFieldValue(mp, "resultScopeCode")
	s.TrackingID = getFieldValue(mp, "trackingId")
	s.Passive = getFieldBooleanValue(mp, "passive")

}
func getFieldValue(mp map[string]interface{}, field string) string {

	if interf, ok := mp[field]; ok {
		if value, ok := interf.(string); ok {
			return value
		}
	}
	return ""
}

func getFieldBooleanValue(mp map[string]interface{}, field string) bool {

	if interf, ok := mp[field]; ok {
		if value, ok := interf.(bool); ok {
			return value
		}
	}
	return false
}
