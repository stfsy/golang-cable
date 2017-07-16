package sortables

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stfsy/golang-assert"
)

const (
	allowAllWildcard  = "/*"
	allowPersWildcard = "/per"
	person            = "/persons/"
	addPerson         = "/persons/add"
	deletePerson      = "/persons/delete"
	deletePerson1     = "/persons/delete/1"
)

func TestSortByLengthDescending(t *testing.T) {
	handlers := SorteableMatchedRequestHandlers{
		MatchedRequestHandler{Match: person},
		MatchedRequestHandler{Match: addPerson}}

	sort.Sort(handlers)

	assert.Equal(t, handlers[0].Match, addPerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 0, addPerson))
	assert.Equal(t, handlers[1].Match, person, fmt.Sprintf("After sorting the handler at idx %d is %s", 1, person))
}

func TestSortByLengthDescending2(t *testing.T) {
	handlers := SorteableMatchedRequestHandlers{
		MatchedRequestHandler{Match: person},
		MatchedRequestHandler{Match: addPerson},
		MatchedRequestHandler{Match: deletePerson}}

	sort.Sort(handlers)

	assert.Equal(t, handlers[0].Match, deletePerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 0, deletePerson1))
	assert.Equal(t, handlers[1].Match, addPerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 1, addPerson))
	assert.Equal(t, handlers[2].Match, person, fmt.Sprintf("After sorting the handler at idx %d is %s", 2, person))
}

func TestSortByLengthDescending3(t *testing.T) {
	handlers := SorteableMatchedRequestHandlers{
		MatchedRequestHandler{Match: person},
		MatchedRequestHandler{Match: addPerson},
		MatchedRequestHandler{Match: deletePerson},
		MatchedRequestHandler{Match: deletePerson1}}

	sort.Sort(handlers)

	assert.Equal(t, handlers[0].Match, deletePerson1, fmt.Sprintf("After sorting the handler at idx %d is %s", 0, deletePerson1))
	assert.Equal(t, handlers[1].Match, deletePerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 1, deletePerson))
	assert.Equal(t, handlers[2].Match, addPerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 2, addPerson))
	assert.Equal(t, handlers[3].Match, person, fmt.Sprintf("After sorting the first handler %d is %s", 3, person))
}

func TestSortByLengthDescending4(t *testing.T) {
	handlers := SorteableMatchedRequestHandlers{
		MatchedRequestHandler{Match: allowAllWildcard},
		MatchedRequestHandler{Match: person},
		MatchedRequestHandler{Match: addPerson},
		MatchedRequestHandler{Match: deletePerson},
		MatchedRequestHandler{Match: deletePerson1}}

	sort.Sort(handlers)

	assert.Equal(t, handlers[0].Match, deletePerson1, fmt.Sprintf("After sorting the handler at idx %d is %s", 0, deletePerson1))
	assert.Equal(t, handlers[1].Match, deletePerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 1, deletePerson))
	assert.Equal(t, handlers[2].Match, addPerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 2, addPerson))
	assert.Equal(t, handlers[3].Match, person, fmt.Sprintf("After sorting the first handler %d is %s", 3, person))
	assert.Equal(t, handlers[4].Match, allowAllWildcard, fmt.Sprintf("After sorting the first handler %d is %s", 4, allowAllWildcard))
}

func TestSortByLengthDescending5(t *testing.T) {
	handlers := SorteableMatchedRequestHandlers{
		MatchedRequestHandler{Match: allowAllWildcard},
		MatchedRequestHandler{Match: allowPersWildcard},
		MatchedRequestHandler{Match: person},
		MatchedRequestHandler{Match: addPerson},
		MatchedRequestHandler{Match: deletePerson},
		MatchedRequestHandler{Match: deletePerson1}}

	sort.Sort(handlers)

	assert.Equal(t, handlers[0].Match, deletePerson1, fmt.Sprintf("After sorting the handler at idx %d is %s", 0, deletePerson1))
	assert.Equal(t, handlers[1].Match, deletePerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 1, deletePerson))
	assert.Equal(t, handlers[2].Match, addPerson, fmt.Sprintf("After sorting the handler at idx %d is %s", 2, addPerson))
	assert.Equal(t, handlers[3].Match, person, fmt.Sprintf("After sorting the first handler %d is %s", 3, person))
	assert.Equal(t, handlers[4].Match, allowPersWildcard, fmt.Sprintf("After sorting the first handler %d is %s", 4, allowPersWildcard))
	assert.Equal(t, handlers[5].Match, allowAllWildcard, fmt.Sprintf("After sorting the first handler %d is %s", 5, allowAllWildcard))
}
