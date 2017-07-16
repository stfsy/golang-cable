package cable

/*
* #  ROADMAP
* ## Content-Type
* 1. Do not ignore params of Content-Type header. E.g. check encoding parameter
* 2. Do not ignore Content-Encoding header. http://greenbytes.de/tech/webdav/rfc2616.html#rfc.section.14.11
* 3. Want to check http://greenbytes.de/tech/webdav/rfc2616.html#rfc.section.14.15 ?
*
* ## Accept
* 1. Accept-Charset
* 2. Accept-Encoding
 */

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/textproto"
	"sort"
	"strings"

	"github.com/stfsy/golang-cable/cable/pluggables"
	"github.com/stfsy/golang-cable/cable/sortables"
)

const (
	defaultContentType = "application/json"
)

var (
	contentTypeHeader = textproto.CanonicalMIMEHeaderKey("content-type")
	acceptHeader      = textproto.CanonicalMIMEHeaderKey("accept")
	plugins           = []pluggables.Plugin{pluggables.JsonPlugin{}, pluggables.XmlPlugin{}}
)

func UnmarshalBody(target interface{}, request *http.Request) (interface{}, error) {
	contentTypes := request.Header[contentTypeHeader]

	/*
	* Any HTTP/1.1 message containing an entity-body SHOULD include a Content-Type header
	* field defining the media type of that body. If and only if the media type is not given
	* by a Content-Type field, the recipient MAY attempt to guess the media type via inspection
	* of its content and/or the name extension(s) of the URI used to identify the resource. If
	* the media type remains unknown, the recipient SHOULD treat it as type "application/octet-stream".
	*
	* https://www.w3.org/Protocols/rfc2616/rfc2616-sec7.html#sec7
	 */
	if len(contentTypes) <= 0 {
		contentTypes = append(contentTypes, "application/octet-stream")
	}

	contentType := contentTypes[0]
	parsedContentType, _, err := mime.ParseMediaType(contentType)

	body, err := ioutil.ReadAll(request.Body)

	if err == nil {

		handled := false

	loop:
		for _, plugin := range plugins {
			for _, consumes := range plugin.Consumes() {
				if parsedContentType == consumes {
					err = plugin.Consume(body, &target)
					handled = true
					break loop
				}
			}
		}

		if !handled {
			msg := fmt.Sprintf("Unmarshalling failed, no marshaller for content type %v", contentType)
			err = errors.New(msg)
			logger.Error(msg)
		}
	}

	return target, err
}

func MarshalBody(target interface{}, request *http.Request) ([]byte, error) {
	/* Accessing the map directy means we have to use the correct casing */
	contentTypes := request.Header[acceptHeader]
	sortableMediaTypes := sortables.SortableMediaTypes{}

	/*
	* If no Accept header field is present, then it is assumed that the client accepts all media types.
	*
	* https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html
	 */
	if len(contentTypes) == 0 {
		contentTypes = append(contentTypes, defaultContentType)
	}

	for _, r := range contentTypes {
		splittedContentTypes := strings.Split(r, ",")

		for _, r := range splittedContentTypes {
			mediatype, params, err := mime.ParseMediaType(r)

			if err == nil {
				sortableMediaTypes = append(sortableMediaTypes, sortables.MediaTypeAndParams{
					Params:    params,
					Mediatype: mediatype})
			}
		}
	}

	sort.Sort(&sortableMediaTypes)

	for _, mediaTypeAndParams := range sortableMediaTypes {
		for _, plugin := range plugins {
			for _, produces := range plugin.Produces() {
				if mediaTypeAndParams.Mediatype == produces {
					return plugin.Produce(target)
				}
			}
		}
	}

	return nil, fmt.Errorf("No Marshaller found for acceptable content types")
}
