package structs

type Configuration struct {
	DBConfig       *DBConf `json:"database,omitempty"`
	JwtKey         string  `json:"jwtkey,omitempty"`
	WebDir         string  `json:"webdir,omitempty"`
	RulebaseDir    string  `json:"rulebasedir,omitempty"`
	Gopathrelative bool    `json:"gopathrelative,omitempty"`
	Loadtestdata   bool    `json:"loadtestdata,omitempty"`
}

type DBConf struct {
	Location string `json:"location"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty`
}

type User struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	Password         string     `json:"password"`
	Firstname        string     `json:"firstname"`
	Lastname         string     `json:"lastname"`
	Locale           string     `json:"locale"`
	AssumptionSet    string     `json:"assumptionSet"`
	Revision         string     `json:"revision"`
	GlobalDictionary Dictionary `json:"globalDictionary"`
	//Documents []string `json:"documents"`
}

type DictionaryEntry struct {
	Value          string `json:"value"`
	Type           string `json:"type"`
	Code           string `json:"code"`
	Category       string `json:"category"`
	DictionaryType string `json:"dictionaryType"`
}

type Dictionary map[string]DictionaryEntry

//Response represents a JSON response from the ducklib server
type Response struct {
	Ok        bool        `json:"ok"`
	Reason    *string     `json:"reason,omitempty"`
	ID        *string     `json:"id,omitempty"`
	Documents *[]Document `json:"documents,omitempty"`
}

type ComplianceResponse struct {
	Ok        bool        `json:"ok"`
	Compliant string      `json:"compliant"`
	Documents []*Document `json:"documents,omitempty"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) FromValueMap(mp map[string]interface{}) {

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
	if assumptionSet, ok := mp["assumptionSet"]; ok {
		u.AssumptionSet = assumptionSet.(string)
	}

	/*	if docs, prs := mp["documents"].([]interface{}); prs {
		u.Documents = make([]string, len(docs))
		for i, v := range docs {
			u.Documents[i] = v.(string)
		}
	}*/
	if dict, prs := mp["dictionary"].(map[string]interface{}); prs {
		u.GlobalDictionary = make(Dictionary)
		u.GlobalDictionary.FromInterfaceMap(dict)
	}

}

func (d Dictionary) FromInterfaceMap(mp map[string]interface{}) {

	//Map looks like
	//dictionary:
	//	map[
	//		microsoft_excel:map[code:microsoft_excel category:1 value:Microsoft Excel type:scope]
	//		microsoft_word:map[category:1 value:Microsoft Word type:scope code:microsoft_word]]

	if d == nil {
		d = make(Dictionary)
	}
	for key, value := range mp {
		var de DictionaryEntry

		value := value.(map[string]interface{})
		if code, ok := value["code"]; ok {
			de.Code = code.(string)
		}
		if tpe, ok := value["type"]; ok {
			de.Type = tpe.(string)
		}
		if val, ok := value["value"]; ok {
			de.Value = val.(string)
		}
		if category, ok := value["category"]; ok {
			de.Category = category.(string)
		}
		if dictionaryType, ok := value["dictionaryType"]; ok {
			de.DictionaryType = dictionaryType.(string)
		}

		d[key] = de
	}

}

type Rulebase struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Revision string `json:"_rev"`
}

//
type Taxonomy map[string][]struct {
	Value    string `json:"value"`
	Code     string `json:"code"`
	Category string `json:"category"`
	Fixed    bool   `json:"fixed"`
}

// HTTPError is an error with an http statuscode, it can also wrap another underlying error
type HTTPError struct {
	Err    string
	Status int
	Cause  error
}

// NewHTTPError returns a httpError which implements the Error interface and has the additional field Status for a http status code.
func NewHTTPError(err string, code int) HTTPError {
	return HTTPError{err, code, nil}
}

func WrapErrWith(err error, herr HTTPError) HTTPError {
	herr.Cause = err
	return herr
}

func (e HTTPError) Error() string {
	return e.Err

}
