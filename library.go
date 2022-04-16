package main

import (
	"library/config"
	"library/server"
	"log"
	"sync"
)

func main() {
	log.Println("Starting Library Server")

	log.Println("Initializing configuration")
	err := config.InitConfig("library", nil)
	if err != nil {
		log.Fatal("Failed to read configuration: %v\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Println("Starting HTTP Server")
		err := server.StartHTTPServer()
		if err != nil {
			log.Fatal("Could not start HTTP Server: %v\n", err)
		}
		log.Println("HTTP Server gracefully terminated")
	}()
	wg.Wait()


	

	log.Println("Library Server stopped")
}