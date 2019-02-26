package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ah8ad3/gateway/pkg/admin"
	"github.com/ah8ad3/gateway/pkg/integrate"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"github.com/ah8ad3/gateway/plugins/ip"

	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/plugins/auth"
	"github.com/ah8ad3/gateway/plugins/ratelimitter"
	"github.com/go-chi/chi"
)

func init() {
	go ratelimitter.CleanupVisitors()
	go ip.UpdateBlockList()
}

// V1 Route function for first method of routing
func V1() *chi.Mux {
	r := chi.NewRouter()

	// Rate limiter per user
	//r.Use(ratelimitter.TestMiddle(10))

	// Ip block Middleware
	r.Use(ip.Middleware(nil))

	r.Get("/", admin.Welcome)

	r.Route("/api/v10/admin", func(r chi.Router) {
		r.Get("/service", admin.GETService)
		r.Get("/service/{service_name}", admin.GETServiceSlug)
		r.Post("/service/{service_name}/{version}/plugin", admin.AddPlugin)
		r.Delete("/service/{service_name}/{version}/plugin", admin.DeletePlugin)
		r.Post("/service", admin.PostService)
		r.Delete("/service", admin.DeleteService)
		r.Put("/service", admin.UpdateService)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", auth.RegisterUser)
		r.Post("/sign", auth.SignJWT)
		r.Post("/check", auth.CheckJwt)
	})

	for _, val := range proxy.Services {
		r.Route(val.Path, func(r chi.Router) {
			for _, plug := range val.Plugins {

				if plug.Active {
					r.Use(plug.Middleware(plug.Config))
				}
			}
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

	for _, val := range integrate.Integrates {
		r.Get(val.Path, func(w http.ResponseWriter, r *http.Request) {
			var result []map[string]interface{}
			_ = result
			for _, val := range integrate.Integrates {
				if val.Path == r.URL.Path {
					for _, service := range val.Join {
						var ser []map[string]interface{}
						_ = ser
						url := r.Host + service
						res, err := integrate.GetIntegrateService(url, r.URL.RawQuery)

						// check if service offline create error cause fixed aggregation
						if err && val.Fixed {
							logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log",
								Description: "One of the services was offline in aggregation"}, RequestURL: r.URL.Path,
								IP: r.RemoteAddr, Time: time.Now()})
							_, _ = w.Write([]byte(`{"error": "Aggregation failed one of the services are offline, log stored"}`))

							return
						}

						_ = json.Unmarshal(res, &ser)
						for _, item := range ser {
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

// V2 is prefix type route
func V2() *chi.Mux {
	r := chi.NewRouter()

	for _, val := range proxy.Services {
		proxy := NewProxy(&val)
		r.HandleFunc(val.Path+"/*", proxy.handleProxy)

	}

	return r
}
