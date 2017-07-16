package cable

import (
	"net/http"
)

func Filter(pattern string, handler func(writer http.ResponseWriter, request *http.Request)) {
	registerFilter(pattern, handler)
}
