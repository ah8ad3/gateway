package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	})
	fmt.Println("Server run at :3000")
	_ = http.ListenAndServe(":3000", r)
}
