package pluggables

type Plugin interface {
	Consumes() []string
	Produces() []string
	Consume(b []byte, i *interface{}) error
	Produce(i interface{}) ([]byte, error)
}
