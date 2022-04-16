package server

import "sync"

type handlerAPI struct {
	requestPool *sync.Pool
	bufferPool *sync.Pool
}