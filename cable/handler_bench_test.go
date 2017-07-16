package cable

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Benchmark200PostBikesSold(b *testing.B) {
	registerGetPersons()
	registerPostPersons()
	registerPostBikes()
	registerPostBikesSold()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		request := httptest.NewRequest(http.MethodPost, "/bikes/sold", strings.NewReader("[]"))
		request.Header.Add("Content-type", "application/json")
		recorder := httptest.NewRecorder()

		HandleRequest(recorder, request)
	}
}
