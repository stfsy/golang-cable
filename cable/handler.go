package cable

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/stfsy/golang-cable/cable/sortables"
	"github.com/stfsy/golang-cable/cable/util/log"
	"go.uber.org/zap"
)

var (
	logger, _ = log.NewLogger()
	handlers  = map[string][]RequestHandler{
		http.MethodDelete: []RequestHandler{},
		http.MethodGet:    []RequestHandler{},
		http.MethodPost:   []RequestHandler{},
		http.MethodPut:    []RequestHandler{},
		http.MethodPatch:  []RequestHandler{}}
	filters = []mappedFilter{}
)

type filter func(writer http.ResponseWriter, request *http.Request)

type mappedFilter struct {
	Pattern *regexp.Regexp
	Handle  filter
}

//
type RequestHandler struct {
	Pattern *regexp.Regexp
	Handle  func(requestEntity RequestEntity, response *ResponseEntity)
}

type ResponseWriter struct {
	http.ResponseWriter
	finished bool
}

func (r *ResponseWriter) WriteHeader(header int) {
	r.ResponseWriter.WriteHeader(header)
	r.finished = true
}

func (r *ResponseWriter) Write(bytes []byte) (int, error) {
	l, err := r.ResponseWriter.Write(bytes)

	if err == nil {
		r.finished = true
	}

	return l, err
}

type RequestEntity struct {
	Request *http.Request
}

type ResponseEntity struct {
	Body    []byte
	Header  map[string]string
	Request *http.Request
	Status  int
}

var (
	wildCardRegex           = regexp.MustCompile(".*")
	four0FourRequestHandler = RequestHandler{
		Pattern: wildCardRegex,
		Handle: func(requestEntity RequestEntity, response *ResponseEntity) {
			response.Status = http.StatusNotFound
		}}

	four0FiveRequestHandler = RequestHandler{
		Pattern: wildCardRegex,
		Handle: func(requestEntity RequestEntity, response *ResponseEntity) {
			response.Status = http.StatusMethodNotAllowed
		}}

	four0SixRequestHandler = RequestHandler{
		Pattern: wildCardRegex,
		Handle: func(requestEntity RequestEntity, response *ResponseEntity) {
			response.Status = http.StatusNotAcceptable
		}}

	four1FiveRequestHandler = RequestHandler{
		Pattern: wildCardRegex,
		Handle: func(requestEntity RequestEntity, response *ResponseEntity) {
			response.Status = http.StatusUnsupportedMediaType
		}}
)

func HandleRequest(writer http.ResponseWriter, request *http.Request) {

	logger.Debug("Handling Request",
		zap.String("Path", request.URL.Path),
		zap.String("Method", request.Method))

	wrappedWriter := ResponseWriter{writer, false}
	handler := findHandler(request.URL.Path, request.Method, writer)
	filter := findFilter(request.URL.Path)

	for _, f := range filter {
		f.Handle(&wrappedWriter, request)

		if wrappedWriter.finished {
			// lets see if its a good idea to return here
			return
		}
	}

	responseEntity := ResponseEntity{Request: request}
	requestEntity := RequestEntity{Request: request}

	handler.Handle(requestEntity, &responseEntity)

	logger.Debug("Handling Response",
		zap.String("Path", request.URL.Path),
		zap.String("Method", request.Method),
		zap.Int("StatusCode", responseEntity.Status),
		zap.String("Response", string(responseEntity.Body)))

	wrappedWriter.WriteHeader(responseEntity.Status)
	wrappedWriter.Write(responseEntity.Body)
}

func findFilter(path string) []mappedFilter {
	matchingfilter := []mappedFilter{}

	for _, r := range filters {
		match := r.Pattern.FindString(path)

		if len(match) > 0 {
			matchingfilter = append(matchingfilter, r)
		}
	}

	return matchingfilter
}

func findHandler(path string, method string, writer http.ResponseWriter) RequestHandler {
	matchingHandlers := findHandlersForPathAndMethod(path, method)

	// got results? call the handler
	// if not we gotta check if the given path has a handler for a different http method
	if len(matchingHandlers) > 0 {
		rh := matchingHandlers[0].RequestHandler.(RequestHandler)
		return rh
	}

	methodsForPath := []string{}

	for _, r := range []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodDelete} {

		if r == method {
			continue
		}

		handlers := findHandlersForPathAndMethod(path, r)

		if len(handlers) > 0 {
			methodsForPath = append(methodsForPath, r)
		}
	}

	if len(methodsForPath) > 0 {
		writer.Header().Add("Allow", strings.Join(methodsForPath, ","))
		return four0FiveRequestHandler
	}

	return four0FourRequestHandler
}

func findHandlersForPathAndMethod(path string, method string) sortables.SorteableMatchedRequestHandlers {
	handlerArray := handlers[method]
	matchingHandlers := sortables.SorteableMatchedRequestHandlers{}

	for _, r := range handlerArray {
		match := r.Pattern.FindString(path)

		if len(match) > 0 {
			matchingHandlers = append(matchingHandlers, sortables.MatchedRequestHandler{
				Match:          match,
				RequestHandler: r})
		}
	}

	sort.Sort(matchingHandlers)

	return matchingHandlers
}

func registerHandler(method string, pattern string, handler func(requestEntity RequestEntity, response *ResponseEntity)) {
	compiledPattern := stringToRegex(pattern)
	handlers[method] = append(handlers[method], RequestHandler{Pattern: compiledPattern, Handle: handler})

	logger.Info("Registered Handler",
		zap.String("Pattern", compiledPattern.String()),
		zap.String("Method", method))
}

func stringToRegex(pattern string) *regexp.Regexp {
	if strings.HasPrefix(pattern, "^") == false {
		pattern = "^" + pattern
	}

	if strings.HasSuffix(pattern, "$") == false {
		// if the pattern is not ending with a wildcard, terminate it
		terminateRegex := true
		if strings.HasSuffix(pattern, "*") {
			terminateRegex = false
		} else if strings.HasSuffix(pattern, "+") {
			terminateRegex = false
		}
		// add an optional trailing slash
		pattern = pattern + "(/)?"

		if terminateRegex {
			pattern = pattern + "$"
		}
	}

	compiledPattern := regexp.MustCompile(pattern)
	compiledPattern.Longest()

	return compiledPattern
}

func registerFilter(pattern string, handler filter) {
	compiledPattern := stringToRegex(pattern)
	filters = append(filters, mappedFilter{Pattern: compiledPattern, Handle: handler})

	logger.Info("Registered Filter",
		zap.String("Pattern", pattern))
}

func reset() {
	handlers = map[string][]RequestHandler{
		http.MethodDelete: []RequestHandler{},
		http.MethodGet:    []RequestHandler{},
		http.MethodPost:   []RequestHandler{},
		http.MethodPut:    []RequestHandler{},
		http.MethodPatch:  []RequestHandler{}}
	filters = []mappedFilter{}
}
