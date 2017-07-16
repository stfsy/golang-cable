package cable

import "net/http"

func Delete(pattern string, handler func(requestEntity RequestEntity, response *ResponseEntity)) {
	registerHandler(http.MethodDelete, pattern, handler)
}

func Get(pattern string, handler func(requestEntity RequestEntity, response *ResponseEntity)) {
	registerHandler(http.MethodGet, pattern, handler)
}

func Post(pattern string, handler func(requestEntity RequestEntity, response *ResponseEntity)) {
	registerHandler(http.MethodPost, pattern, handler)
}

func Put(pattern string, handler func(requestEntity RequestEntity, response *ResponseEntity)) {
	registerHandler(http.MethodPut, pattern, handler)
}

func Patch(pattern string, handler func(requestEntity RequestEntity, response *ResponseEntity)) {
	registerHandler(http.MethodPatch, pattern, handler)
}
