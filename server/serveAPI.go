package server

import (
	"bytes"
	"context"
	"io"
	"library/handler"
	"library/util"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type Request struct {
	Authorization string
	Body          io.Reader
	URL           *url.URL
	Method        string
}


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

func (handlerAPI *handlerAPI) ServerHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	startTime := time.Now()
	ctx := context.Background()

	authorization := httpRequest.Header.Get("Authorization")
	requestBody, logRequestBody := handlerAPI.getRequestBody(httpRequest.Body)

	request := handlerAPI.requestPool.Get().(*handler.Request)
	request.Authorization = authorization
	request.Body = requestBody
	request.URL = httpRequest.URL
	request.Method = httpRequest.Method

	var response interface{}
	var err error

	var httpResponseStatus int
	var logResponseBody string

	defer func() {
		handlerAPI.requestPool.Put(request)

		duration := time.Now().Sub(startTime)
		
		log.Printf("%v: status=%v method=%v uri=%v duration=%v", startTime, httpResponseStatus, 
			httpRequest.Method, httpRequest.RequestURI, duration)
		log.Printf("request=%v", logRequestBody)

		if response != nil {
			log.Printf("response=%v", logResponseBody)
		}

		if err == nil {
			httpResponseStatus = http.StatusOK
		} else {
			isError, errorCode, cause, errorType := util.IsError(err)

			if isError == true {
				response = util.ErrorResponse{
					ErrorCode: errorCode,
					Cause: cause,
				}

				httpResponseStatus = util.MapErrorTypeToHTTPStatus(errorType)
			} else {
				response = nil
				httpResponseStatus = util.MapErrorTypeToHTTPStatus(err)
			}
		}

		responseBuffer := handlerAPI.bufferPool.Get(*bytes.Buffer)
		responseBuffer.Reset()

		response != nil {
			shouldGZip := err == nil
			logRequestBody = handlerAPI.makeResponseBody(shouldGZip, responseBuffer, response)

			if shouldGZip == true {
				writer.Header().Set("Content-Encoding", "gzip")
			}

			writer.Header().Set("Content-Type", "application/json; charset=utf-8")

			writer.Header().Set("Content-Length", strconv.Itoa(responseBuffer.Len()))

		}

		writer.WriteHeader(httpResponseStatus)
		writer.Write(responseBuffer.Bytes())
	   
		handlerAPI.bufferPool.Put(responseBuffer)

	}()
	response, err = handler.Handle(ctx, request)


}

func (handlerAPI *handlerAPI) getRequestBody(reader io.Reader) (io.Reader, string) {
	buffer := handlerAPI.bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()

	io.Copy(buffer, reader)
	body := buffer.String()

	handlerAPI.bufferPool.Put(buffer)

	return strings.NewReader(body), body
}


func (handlerAPI *handlerAPI) makeResponseBody(shouldGZip bool, writer io.Writer, response interface{}) string {
 
	if response == nil {
	   return ""
	}
	
	respRawBody := handlerAPI.bufferPool.Get().(*bytes.Buffer)
	respRawBody.Reset()
 
	if shouldGZip == true {
	   gzipper := gzip.NewWriter(writer)
	   writeResponse(io.MultiWriter(respRawBody, gzipper), response)
 
	   gzipper.Close()
	} else {
	   writeResponse(io.MultiWriter(respRawBody, writer), response)
	}
 
	rawBody := trimEOL(respRawBody.String())
 
	handlerAPI.bufferPool.Put(respRawBody)
 
	return rawBody
  }

  func trimEOL(json string) string {
	n := len(json)

	if n > 0 && json[n-1] == '\n' {
	   return json[:n-1]
	}

	return json
  }

  func writeResponse(writer io.Writer, response interface{}) {

	switch resp := response.(type) {
	case []byte:
	   writer.Write(resp)
	default:
	   json.NewEncoder(writer).Encode(response)
	}
}
