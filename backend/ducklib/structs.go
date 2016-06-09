package ducklib

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Locale    string   `json:"locale"`
	Revision  string   `json:"_rev"`
	Documents []string `json:"documents"`
}

type Response struct {
	Ok     bool    `json:"ok"`
	Reason *string `json:"reason,omitempty"`
	ID     *string `json:"id,omitempty"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) fromValueMap(mp map[string]interface{}) {

	if id, ok := mp["_id"]; ok {
		u.ID = id.(string)
	}
	if rev, ok := mp["_rev"]; ok {
		u.Revision = rev.(string)
	}
	if name, ok := mp["firstname"]; ok {
		u.Firstname = name.(string)
	}
	if owner, ok := mp["lastname"]; ok {
		u.Lastname = owner.(string)
	}
	if owner, ok := mp["password"]; ok {
		u.Password = owner.(string)
	}
	if owner, ok := mp["email"]; ok {
		u.Email = owner.(string)
	}
	if locale, ok := mp["locale"]; ok {
		u.Locale = locale.(string)
	}

	if docs, prs := mp["documents"].([]interface{}); prs {
		u.Documents = make([]string, len(docs))
		for i, v := range docs {
			u.Documents[i] = v.(string)
		}
	}

}

type Rulebase struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Revision string `json:"_rev"`
}

type Document struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Revision   string      `json:"_rev"`
	Owner      string      `json:"owner"`
	Locale     string      `json:"locale"`
	Statements []Statement `json:"statements"`
}

type Statement struct {
	UseScope     string `json:"useScope"`
	Qualifier    string `json:"qualifier"`
	DataCategory string `json:"dataCategory"`
	SourceScope  string `json:"sourceScope"`
	Action       string `json:"action"`
	ResultScope  string `json:"resultScope"`
	TrackingID   string `json:"trackingId"`
	Passive      bool   `json:"passive"`
}

func (d *Document) fromValueMap(mp map[string]interface{}) {

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

	if stmts, prs := mp["statements"].([]interface{}); prs {
		d.Statements = make([]Statement, len(stmts))
		for i, stmt := range stmts {
			s := new(Statement)
			s.fromInterfaceMap(stmt.(map[string]interface{}))
			d.Statements[i] = *s
		}
	}

}

func (s *Statement) fromInterfaceMap(mp map[string]interface{}) {

	s.UseScope = getFieldValue(mp, "useScope")
	s.Qualifier = getFieldValue(mp, "qualifier")
	s.DataCategory = getFieldValue(mp, "dataCategory")
	s.SourceScope = getFieldValue(mp, "sourceScope")
	s.Action = getFieldValue(mp, "action")
	s.ResultScope = getFieldValue(mp, "resultScope")
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
