package multiplexer

import "net/http"

type Multiplexer = http.ServeMux

func New() *Multiplexer {
	return new(Multiplexer)
}
