package sortables

import (
	"sort"
	"testing"

	assert "github.com/stfsy/golang-assert"
)

var acceptAll = MediaTypeAndParams{Mediatype: "*/*;", Params: map[string]string{"q": "0.8"}}
var acceptPng = MediaTypeAndParams{Mediatype: "image/apng", Params: map[string]string{"q": "0.9"}}
var acceptXML = MediaTypeAndParams{Mediatype: "application/xml", Params: map[string]string{}}

func TestDefaultQualityIs1(t *testing.T) {
	s := SortableMediaTypes{acceptAll, acceptPng, acceptXML}

	sort.Sort(s)

	assert.Equal(t, s[0].Mediatype, acceptXML.Mediatype, "Mediatype XML should be at index 0")
	assert.Equal(t, s[1].Mediatype, acceptPng.Mediatype, "Mediatype PNG should be at index 1")
	assert.Equal(t, s[2].Mediatype, acceptAll.Mediatype, "Wildcard Mediatype should be at index 2")
}

func TestSortByQualityAscending(t *testing.T) {
	s := SortableMediaTypes{acceptAll, acceptPng}

	sort.Sort(s)

	assert.Equal(t, s[0].Mediatype, acceptPng.Mediatype, "Mediatype PNG should be at index 0")
	assert.Equal(t, s[1].Mediatype, acceptAll.Mediatype, "Wildcard Mediatype should be at index 1")
}
