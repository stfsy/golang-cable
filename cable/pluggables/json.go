package pluggables

import (
	"encoding/json"
)

type JsonPlugin struct {
	Plugin
}

var (
	jsonConsumes = []string{"*/*", "application/*", "application/json"}
	jsonProduces = []string{"*/*", "application/*", "application/json"}
)

func (j JsonPlugin) Consumes() []string { return jsonConsumes }
func (j JsonPlugin) Produces() []string { return jsonProduces }

func (j JsonPlugin) Consume(b []byte, i *interface{}) error {
	return json.Unmarshal(b, i)
}

func (j JsonPlugin) Produce(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}
