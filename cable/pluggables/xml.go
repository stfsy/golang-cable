package pluggables

import (
	"encoding/xml"
)

type XmlPlugin struct {
	Plugin
}

var (
	xmlConsumes = []string{"application/xml"}
	xmlProduces = []string{"application/xml"}
)

func (x XmlPlugin) Consumes() []string { return xmlConsumes }
func (x XmlPlugin) Produces() []string { return xmlProduces }

func (x XmlPlugin) Consume(b []byte, i *interface{}) error {
	return xml.Unmarshal(b, i)
}

func (x XmlPlugin) Produce(i interface{}) ([]byte, error) {
	return xml.Marshal(i)
}
