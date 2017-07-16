package cable

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	assert "github.com/stfsy/golang-assert"
)

/*
* # ROADMAP
* - Test returns 415 if content type cannot be handled
* - Test returns 406 if accepts type cannot be handled
* - Test returns 500 if panics
* - Test add global http headers
 */

func Test404NoHandlerRegistered(t *testing.T) {
	reset()
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 404, "Response status code is 404")
}

func Test405NoHandlerRegisteredForMethod(t *testing.T) {
	reset()
	registerGetPersons()
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 405, "Response status code is 405")
}

func Test405AndAllowHeaderNoHandlerRegisteredForMethod(t *testing.T) {
	reset()
	registerGetPersons()
	registerDeletePersons()
	registerPatchPersons()
	registerPutPersons()
	request := httptest.NewRequest(http.MethodConnect, "/persons", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 405, "Response status code is 405")

	registeredMethods := []string{http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete}
	allowedHeaders := recorder.Header().Get("Allow")

loop:
	for _, allowedHeader := range strings.Split(allowedHeaders, ",") {
		for _, registeredMethod := range registeredMethods {
			if allowedHeader == registeredMethod {
				continue loop
			}
		}

		t.Errorf("Expected %v to contain any of %v", allowedHeader, registeredMethods)
	}
}

func Test200PostPersonsNoTrailingSlash(t *testing.T) {
	reset()
	registerGetPersons()
	registerPostPersons()
	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Some Persons", "Response body includes Some Persons")
}

func Test200PostPersonsTrailingSlash(t *testing.T) {
	reset()
	registerGetPersons()
	registerPostPersons()
	request := httptest.NewRequest(http.MethodPost, "/persons/", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Some Persons", "Response body includes Some Persons")
}

func Test200PostBikes(t *testing.T) {
	reset()
	registerGetPersons()
	registerPostPersons()
	registerPostBikes()

	request := httptest.NewRequest(http.MethodPost, "/bikes/", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Bikes", "Response body includes Bikes")
}

func Test200PostBikesSold(t *testing.T) {
	reset()
	registerGetPersons()
	registerPostPersons()
	registerPostBikes()
	registerPostBikesSold()

	request := httptest.NewRequest(http.MethodPost, "/bikes/sold", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Sold!", "Response body includes Sold!")
}

func Test200PutPersons(t *testing.T) {
	reset()
	registerPostPersons()
	registerGetPersons()
	registerPutPersons()
	registerPatchPersons()
	registerDeletePersons()

	request := httptest.NewRequest(http.MethodPut, "/persons", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Put Some Persons", "Response body includes Put")
}

func Test200PatchPersons(t *testing.T) {
	reset()
	registerPostPersons()
	registerGetPersons()
	registerPutPersons()
	registerPatchPersons()
	registerDeletePersons()

	request := httptest.NewRequest(http.MethodPatch, "/persons", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Patched Some Persons", "Response body includes Patched")
}

func Test200DeletePersons(t *testing.T) {
	reset()
	registerPostPersons()
	registerGetPersons()
	registerPutPersons()
	registerPatchPersons()
	registerDeletePersons()

	request := httptest.NewRequest(http.MethodDelete, "/persons", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "Deleted Some Persons", "Response body includes Deleted")
}

func Test200WWildcard1(t *testing.T) {
	reset()
	registerPostCWildcard()

	request := httptest.NewRequest(http.MethodPost, "/ca", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "C!", "Response body includes Deleted")
}

func Test200WWildcard2(t *testing.T) {
	reset()
	registerPostCWildcard()

	request := httptest.NewRequest(http.MethodPost, "/ca/cb/1234", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "C!", "Response body includes Deleted")
}

func Test200WWildcard3(t *testing.T) {
	reset()
	registerPostDWildcard()

	request := httptest.NewRequest(http.MethodPost, "/d", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "D!", "Response body includes Deleted")
}

func Test200WWildcard4(t *testing.T) {
	reset()
	registerPostDWildcard()

	request := httptest.NewRequest(http.MethodPost, "/deeeab1", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "D!", "Response body includes Deleted")
}

func Test200WWildcard5(t *testing.T) {
	reset()
	registerPostDWildcard()

	request := httptest.NewRequest(http.MethodPost, "/deeeab1/bcabc", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 200, "Response status code is 200")
	assert.Equal(t, recorder.Body.String(), "D!", "Response body includes Deleted")
}

func Test200WWildcard6(t *testing.T) {
	reset()
	registerPostCWildcard()

	request := httptest.NewRequest(http.MethodPost, "/def", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 404, "Response status code is 404")
}

func Test200WWildcard7(t *testing.T) {
	reset()
	registerPostDWildcard()

	request := httptest.NewRequest(http.MethodPost, "/Leeeab1/bcabc", strings.NewReader("[]"))
	request.Header.Add("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 404, "Response status code is 404")
}

func registerGetPersons() {
	Get("/persons", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
	})
}

func registerPostPersons() {
	Post("/persons", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("Some Persons")
	})
}

func registerPutPersons() {
	Put("/persons", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("Put Some Persons")
	})
}

func registerPatchPersons() {
	Patch("/persons", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("Patched Some Persons")
	})
}

func registerDeletePersons() {
	Delete("/persons", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("Deleted Some Persons")
	})
}

func registerPostBikes() {
	Post("/bikes/*", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("Bikes")
	})
}

func registerPostBikesSold() {
	Post("/bikes/sold", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("Sold!")
	})
}

func registerPostCWildcard() {
	/* /c* matches everything */
	Post("/c.*", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("C!")
	})
}

func registerPostDWildcard() {
	Post("/d+", func(req RequestEntity, resp *ResponseEntity) {
		resp.Status = 200
		resp.Body = []byte("D!")
	})
}
