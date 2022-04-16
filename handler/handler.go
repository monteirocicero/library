package handler

import (
	"io"
	"net/url"
)

type Request struct {
	Authorization string
	Body io.Reader
	URL *url.URL
	Method string
}