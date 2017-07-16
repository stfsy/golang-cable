package sortables

type SorteableMatchedRequestHandlers []MatchedRequestHandler

type MatchedRequestHandler struct {
	// using interface{} as request handlers type
	// to keep this type  and the request handler itself loosely coupled
	Match          string
	RequestHandler interface{}
}

func (s SorteableMatchedRequestHandlers) Len() int {
	return len(s)
}

func (s SorteableMatchedRequestHandlers) Less(i, j int) bool {
	return len(s[i].Match) > len(s[j].Match)
}

func (s SorteableMatchedRequestHandlers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
