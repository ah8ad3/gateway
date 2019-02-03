package routes

import (
	"encoding/json"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/integrate"
	"github.com/ah8ad3/gateway/plugins/ip"
	"net/http"
	"strings"
	"time"

	"github.com/ah8ad3/gateway/pkg/auth"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/plugins/ratelimitter"
	"github.com/go-chi/chi"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _ = w.Write(getServices())
}

func init() {
	go ratelimitter.CleanupVisitors()
	go ip.UpdateBlockList()
}

// V1 Route function for first method of routing
func V1() *chi.Mux {
	r := chi.NewRouter()

	// Rate limiter per user
	r.Use(ratelimitter.LimitMiddleware)

	// Ip block Middleware
	r.Use(ip.InfoMiddleware)

	r.Get("/", welcome)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", auth.RegisterUser)
		r.Post("/sign", auth.SignJWT)
		r.Post("/check", auth.CheckJwt)
	})

	for _, val := range Services {
		r.Route(val.Path, func(r chi.Router) {
			for _, url := range val.Urls {
				switch url.Method {
				case "GET":
					r.Get(url.Path, func(writer http.ResponseWriter, request *http.Request) {

						// remove path form url and send to service and serve answer
						splitRoute := strings.Split(request.URL.Path, "/")
						route := strings.Join(splitRoute[2:], "/")
						if route == "" {
							route = "/"
						} else {
							route = "/" + route
						}

						logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log"}, RequestURL: request.URL.Path,
							IP: request.RemoteAddr, Time: time.Now()})

						writer.Header().Set("Content-Type", "application/json")


						server := findService(splitRoute[1])
						body, code := GetService(server, route, request.URL.RawQuery)
						writer.WriteHeader(code)
						_, _ = writer.Write(body)
					})

				case "POST":
					r.Post(url.Path, func(writer http.ResponseWriter, request *http.Request) {
						// remove path form url and send to service and serve answer
						splitRoute := strings.Split(request.URL.Path, "/")[2:]
						route := strings.Join(splitRoute, "/")
						if route == "" {
							route = "/"
						} else {
							route = "/" + route
						}
						_ = request.ParseForm()

						m := make(map[string]interface{})
						for key, value := range request.Form {
							m[key] = strings.Join(value, "")
						}

						data, _ := json.Marshal(m)

						logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log"}, RequestURL: request.URL.Path,
							IP: request.RemoteAddr, Time: time.Now()})

						writer.Header().Set("Content-Type", "application/json")
						server := findService(splitRoute[1])
						body, code := PostService(server, route, data)
						writer.WriteHeader(code)
						_, _ = writer.Write(body)
					})

				case "PUT":
					// not implemented now
					r.Put(url.Path, func(writer http.ResponseWriter, request *http.Request) {
						_, _ = writer.Write([]byte("hello"))
					})

				case "DELETE":
					// not implemented now
					r.Delete(url.Path, func(writer http.ResponseWriter, request *http.Request) {
						_, _ = writer.Write([]byte("hello"))
					})

				default:
					logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical",
						Description: fmt.Sprintf("Bad url method in service %s", val.Name)},
						Pkg: "auth", Time: time.Now()})
					//log.Fatal("Bad url method in service ", val.Name)

				}

			}
		})
	}

	for _, val := range integrate.Integrates{
		r.Get(val.Path, func(w http.ResponseWriter, r *http.Request) {
			var result []map[string]interface{}
			_ = result
			for _, val := range integrate.Integrates{
				if val.Path == r.URL.Path {
					for _, service := range val.Join {
						var ser []map[string]interface{}
						_ = ser
						url := r.Host + service
						res, err :=GetIntegrateService(url, r.URL.RawQuery)

						// check if service offline create error cause fixed aggregation
						if err && val.Fixed {
							logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log",
								Description: "One of the services was offline in aggregation"}, RequestURL: r.URL.Path,
								IP: r.RemoteAddr, Time: time.Now()})
							_, _ = w.Write([]byte(`{"error": "Aggregation failed one of the services are offline, log stored"}`))

							return
						}

						_ = json.Unmarshal(res, &ser)
						for _, item := range ser{
							result = append(result, item)
						}
					}
				}
			}
			// check if all the sevices are offline
			if result == nil {
				logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log",
					Description: "all of the services was offline in aggregation"}, RequestURL: r.URL.Path,
					IP: r.RemoteAddr, Time: time.Now()})
				_, _ = w.Write([]byte(`{"error": "Aggregation failed all of the services are offline, log stored"}`))

				return
			}
			jData, _ := json.Marshal(result)
			_, _ = w.Write(jData)
			return
		})
	}
	return r
}
