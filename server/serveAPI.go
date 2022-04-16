package server

import (
	"bytes"
	"library/handler"
	"sync"
)

func newHandlerAPI() (api *handlerAPI) {
	api = &handlerAPI{}

	api.requestPool = &sync.Pool{
		New: func() interface{} {
			return new(handler.Request)
		},
	}

	api.bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	return
}