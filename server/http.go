package server

import (
	"context"
	"library/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	StartHTTPServer = startHTTPServer
)

func startHTTPServer() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/api/", newHandlerAPI())

	server := http.Server{}
	server.ReadHeaderTimeout = config.GetHTTPReadTimeout()
	server.WriteTimeout = config.GetHTTPWriteTimeout()
	server.Addr = config.GetHTTPServerAddress()
	server.Handler = mux

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error shuting down. %v\n", err)
		}
	}()

	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}

	return
}