package routes

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strings"
)

func welcome (w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("welcome"))
}

func Routes() *chi.Mux{
	r := chi.NewRouter()
	r.Get("/",  welcome)

	for _, val := range Service{
		r.Route(val.Path, func(r chi.Router) {
			for _, url := range val.Urls {
				switch url.Method {
				case "GET":
					r.Get(url.Path, func(writer http.ResponseWriter, request *http.Request) {

						// remove path form url and send to service and serve answer
						splitRoute := strings.Split(request.URL.Path, "/")[2:]
						route := strings.Join(splitRoute, "/")
						if route == "" {
							route = "/"
						}

						writer.Header().Set("Content-Type", "application/json")
						body := GetService(val.Server, route, url.Method, request.URL.RawQuery)
						_, _ = writer.Write(body)
					})

				case "POST":
					r.Post(url.Path, func(writer http.ResponseWriter, request *http.Request) {
						// remove path form url and send to service and serve answer
						splitRoute := strings.Split(request.URL.Path, "/")[2:]
						route := strings.Join(splitRoute, "/")
						if route == "" {
							route = "/"
						}
						_ = request.ParseForm()


						for key, value := range request.Form {
							fmt.Println(key, value)
						}

						_, _ = writer.Write([]byte("hello"))
					})

				case "PUT":
					r.Put(url.Path, func(writer http.ResponseWriter, request *http.Request) {
						_, _ = writer.Write([]byte("hello"))
					})

				case "DELETE":
					r.Delete(url.Path, func(writer http.ResponseWriter, request *http.Request) {
						_, _ = writer.Write([]byte("hello"))
					})

				default:
					log.Fatal("Bad url method in service ", val.Name)

				}

			}
		})
	}
	return r
}
