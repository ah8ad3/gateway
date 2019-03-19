package routes

import (
	"fmt"
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
	//r.Use(ip.Middleware(nil))

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
					r.Get(url.Path, getProxyHttp)

				case "POST":
					r.Post(url.Path, postProxyHttp)

				case "PUT":
					// not implemented now
					r.Put(url.Path, putProxyHttp)

				case "DELETE":
					// not implemented now
					r.Delete(url.Path, deleteProxyHttp)

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
		r.Get(val.Path, integrateProxyHttp)
	}
	return r
}

// V2 is prefix type route
func V2() *chi.Mux {
	r := chi.NewRouter()

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
		_proxy := NewProxy(val)
		r.Route(_proxy.service.Path, func(r chi.Router) {
			for _, plug := range val.Plugins{
				r.Use(plug.Middleware(plug.Config))
			}
			r.HandleFunc("/*", _proxy.handleProxy)
		})
	}
	return r
}
