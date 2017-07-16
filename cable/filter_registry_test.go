package cable

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	assert "github.com/stfsy/golang-assert"
)

func TestNoHandlerNoFilterResponds404(t *testing.T) {
	reset()

	request := httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 404, "Response status code is 201")
}

func TestNoHandlerFilterResponds(t *testing.T) {
	reset()
	registerRegisterFilter()

	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 201, "Response status code is 201")
}

func TestNoHandlerFilterResponds2(t *testing.T) {
	reset()
	registerLoginFilter()

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 403, "Response status code is 403")
}

func TestFilterExecutedInOrderAndReturnEarly(t *testing.T) {
	reset()
	registerRegisterFilter()
	registerRegister500Filter()

	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("[]"))
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, request)

	assert.Equal(t, recorder.Code, 201, "Response status code is 201")
}

func registerRegisterFilter() {
	Filter("/register", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(201)
	})
}

func registerRegister500Filter() {
	Filter("/register", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(500)
	})
}

func registerLoginFilter() {
	Filter("/login", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
	})
}
