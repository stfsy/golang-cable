package cable

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stfsy/golang-assert"
)

type Person struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
}

var jsonBody = "{\"firstName\":\"Mario\",\"secondName\":\"Micelli\"}"
var emptyJSONList = "[]"
var unexpectedBody = "{ \"street\": \"nowhere\", \"zip\": \"12345\"}"
var xmlBody = "<Person><FirstName>Mario</FirstName><SecondName>Micelli</SecondName></Person>"

func TestUnmarshallJson(t *testing.T) {
	person, err := Unmarshall("application/json", jsonBody)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "Mario", "Name sould be Mario")
	assert.Equal(t, person.SecondName, "Micelli", "Second name should be Micelli")
}

func TestMarshallJson(t *testing.T) {
	stringBody, err := Marshall("application/json")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, jsonBody, "String body should be JSON")
}

func TestMarshallEmptyJsonList(t *testing.T) {
	stringBody, err := MarshallEmptyList("application/json")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, emptyJSONList, "String Body should an empty JSON List")
}

func TestUnmarshallUnexpectedJson(t *testing.T) {
	person, err := UnmarshallUnexpectedBody("application/json")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "", "First name should be empty")
	assert.Equal(t, person.SecondName, "", "Second name should be empty")
}

func TestUnmarshallJsonWithParameters(t *testing.T) {
	person, err := Unmarshall("application/json; charset=ISO-8859-4", jsonBody)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "Mario", "Name sould be Mario")
	assert.Equal(t, person.SecondName, "Micelli", "Second name should be Micelli")
}

func TestUnmarshallJsonWithWildcardContentType(t *testing.T) {
	person, err := Unmarshall("*/*", jsonBody)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "Mario", "Name sould be Mario")
	assert.Equal(t, person.SecondName, "Micelli", "Second name should be Micelli")
}

func TestUnmarshallJsonWithApplicationWildcardContentType(t *testing.T) {
	person, err := Unmarshall("application/*", jsonBody)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "Mario", "Name sould be Mario")
	assert.Equal(t, person.SecondName, "Micelli", "Second name should be Micelli")
}

func TestMarshallJsonApplicationWildcard(t *testing.T) {
	stringBody, err := Marshall("application/xml;q=0.8", "application/*")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, jsonBody, "String body should be JSON")
}

func TestMarshallJsonWildcard(t *testing.T) {
	stringBody, err := Marshall("application/xml;q=0.8", "*/*")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, jsonBody, "String body should be JSON")
}

func TestMarshallJsonQDefaultParameter(t *testing.T) {
	stringBody, err := Marshall("application/xml;q=0.8", "application/json")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, jsonBody, "String body should be JSON")
}

func TestMarshallJsonQParameter(t *testing.T) {
	stringBody, err := Marshall("application/xml;q=0.8", "application/json;q=0.9")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, jsonBody, "String body should be JSON")
}

func TestUnmarshallJsonWithWildcardApplicationContentType(t *testing.T) {
	person, err := Unmarshall("application/*", jsonBody)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "Mario", "Name sould be Mario")
	assert.Equal(t, person.SecondName, "Micelli", "Second name should be Micelli")
}

func TestUnmarshallJsonNoAcceptHeader(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader(jsonBody))

	bytes, err := MarshalBody(Person{FirstName: "Mario", SecondName: "Micelli"}, request)
	stringBody := string(bytes)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, jsonBody, "String body should be JSON")
}

func TestUnmarshallXml(t *testing.T) {
	person, err := Unmarshall("application/xml", xmlBody)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, person.FirstName, "Mario", "Name sould be Mario")
	assert.Equal(t, person.SecondName, "Micelli", "Second name should be Micelli")
}

func TestUnmarshallUnknownContentType(t *testing.T) {
	_, err := Unmarshall("got/bb", xmlBody)

	assert.NotEqual(t, err, nil, "Error should  not be nil")
}

func TestMarshallXml(t *testing.T) {
	stringBody, err := Marshall("application/xml")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, xmlBody, "String Body should be xml")
}

func TestMarshallXmlMultipleAcceptTypes(t *testing.T) {
	stringBody, err := Marshall("text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/webp", "image/apng", "*/*;q=0.8")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, xmlBody, "String Body should be xml")
}

func TestMarshallXmlMultipleAcceptTypesInHeaderField(t *testing.T) {
	stringBody, err := Marshall("text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, xmlBody, "String Body should be xml")
}

func TestMarshallXmlNotJsonQParameter(t *testing.T) {
	stringBody, err := Marshall("application/json;q=0.1", "application/xml;q=0.9")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, xmlBody, "String Body should be xml")
}

func TestMarshallXmlNotJsonQDefaultParameter(t *testing.T) {
	stringBody, err := Marshall("application/json;q=0.1", "application/xml")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, xmlBody, "String Body should be xml")
}

func TestMarshallMultipleXmlAcceptTypesReversed(t *testing.T) {
	stringBody, err := Marshall("*/*;q=0.8", "image/apng", "image/webp", "application/xml;q=0.9", "application/xhtml+xml", "text/html")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, stringBody, xmlBody, "String Body should be xml")
}

func TestMarshallErrUnknownMediatype(t *testing.T) {
	_, err := Marshall("xyz/abc")

	assert.NotEqual(t, err, nil, "Error should not be nil")
}

func Unmarshall(contentType string, content string) (*Person, error) {
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader(content))
	request.Header.Add("content-type", contentType)

	p, err := UnmarshalBody(&Person{}, request)
	person := p.(*Person)

	return person, err
}

func UnmarshallUnexpectedBody(contentType string) (*Person, error) {
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader(unexpectedBody))
	request.Header.Add("content-type", "application/json")

	p, err := UnmarshalBody(&Person{}, request)
	person := p.(*Person)

	return person, err
}

func Marshall(accepts ...string) (string, error) {
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader(jsonBody))

	for _, r := range accepts {
		request.Header.Add("accept", r)
	}

	bytes, err := MarshalBody(Person{FirstName: "Mario", SecondName: "Micelli"}, request)
	stringBody := string(bytes)

	return stringBody, err
}

func MarshallEmptyList(accepts string) (string, error) {
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader(emptyJSONList))
	request.Header.Add("accept", accepts)

	bytes, err := MarshalBody([]Person{}, request)
	stringBody := string(bytes)

	return stringBody, err
}
