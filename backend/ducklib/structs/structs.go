package structs

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

//Document struct contains among other things a set of statements, a Unique ID, a name and an owner
//The Owner field is a foreign key to a User.ID
type Document struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Revision      string      `json:"revision"`
	Owner         string      `json:"owner"`
	Locale        string      `json:"locale"`
	Description   string      `json:"description"`
	AssumptionSet string      `json:"assumptionSet"`
	Statements    []Statement `json:"statements"`
	Dictionary    Dictionary  `json:"dictionary"`
}

//A Statement struct represents one Statement in a document
type Statement struct {
	UseScopeCode     string  `json:"useScopeCode"`
	QualifierCode    string  `json:"qualifierCode"`
	DataCategoryCode string  `json:"dataCategoryCode"`
	SourceScopeCode  string  `json:"sourceScopeCode"`
	ActionCode       string  `json:"actionCode"`
	ResultScopeCode  string  `json:"resultScopeCode"`
	TrackingID       string  `json:"trackingId"`
	Tag              *string `json:"tag,omitempty`
	Passive          bool    `json:"passive"`
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
	Compliant   string      `json:"compliant"`
	Explanation interface{} `json:"explanation"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
