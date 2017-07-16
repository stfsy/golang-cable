package sortables

import "strconv"

type SortableMediaTypes []MediaTypeAndParams

func (s SortableMediaTypes) Len() int {
	return len(s)
}

func (s SortableMediaTypes) Less(i, j int) bool {
	iParams := s[i].Params
	jParams := s[j].Params

	quality := "q"

	/*
	* Each media-range MAY be followed by one or more accept-params,
	* beginning with the "q" parameter for indicating a relative quality factor.
	* The first "q" parameter (if any) separates the media-range parameter(s) from
	* the accept-params. Quality factors allow the user or user agent to indicate the
	*  relative degree of preference for that media-range, using the qvalue scale from
	* 0 to 1 (section 3.9). The default value is q=1.
	*
	* https: //www.w3.org/Protocols/rfc2616/rfc2616-sec14.html
	 */
	iQuality := 1.0
	jQuality := 1.0

	if q := iParams[quality]; len(q) == 3 {
		f, err := strconv.ParseFloat(q, 32)
		if err == nil {
			iQuality = f
		}
	}

	if q := jParams[quality]; len(q) == 3 {
		f, err := strconv.ParseFloat(q, 32)
		if err == nil {
			jQuality = f
		}
	}

	return iQuality > jQuality
}

func (s SortableMediaTypes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type MediaTypeAndParams struct {
	Mediatype string
	Params    map[string]string
}
