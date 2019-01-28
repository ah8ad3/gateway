package routes

import (
	"github.com/go-chi/chi"
	"net/http"
)

func welcome (w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("welcome"))
}

func Routes() *chi.Mux{
	r := chi.NewRouter()
	r.Get("/",  welcome)
	return r
}
