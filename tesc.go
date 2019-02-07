package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func m2() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":2017"
	}

	h := &http.Server{Addr: addr, Handler: &server{}}

	logger := log.New(os.Stdout, "", 0)

	go func() {
		logger.Printf("Listening on http://0.0.0.0%s\n", addr)

		if err := h.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	h.Shutdown(context.Background())

	<-stop

	logger.Println("\nShutting down the server...")

	h.Shutdown(context.Background())

	logger.Println("Server gracefully stopped")
}

func main1() {
	m2()
}

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}