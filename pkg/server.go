package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/pkg/integrate"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"github.com/ah8ad3/gateway/pkg/routes"
	"github.com/ah8ad3/gateway/plugins"
	"github.com/ah8ad3/gateway/plugins/auth"
	"github.com/go-chi/chi"
)

var serLocation string

func init() {
	if os.Getenv("TEST") == "1" {
		serLocation = "./../services.json"
	} else {
		serLocation = "services.json"
	}
}

func settings() {
	logger.OpenConnection()
	// Require for auth Database staff
	auth.OpenAuthCollection()

	// register all plugins to gateway
	plugins.RegisterPlugins()

	proxy.LoadServices(false, serLocation)
	proxy.CheckServices(false)

	integrate.LoadIntegration()

	// check all service available every one hour
	go proxy.HealthCheck()
}

// RUN for run server
func RUN(ip string, port string, route string) {
	str, _ := db.LoadSecretKey()
	db.SecretKey = str
	if db.SecretKey == "" {
		log.Fatal("Secret Key not Found, generate it by\n  \t gateway secret")
	}
	var r *chi.Mux
	settings()
	if route == "v1" {
		r = routes.V1()
	} else {
		r = routes.V2()
	}

	if port == "" {
		port = "3000"
	}
	if ip == "" {
		ip = "0.0.0.0"
	}

	listen := ip + ":" + port

	hs := setup(listen, r)

	fmt.Println(fmt.Sprintf("Listening on http://%s\n", hs.Addr))

	err := hs.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}

}

func setup(url string, r *chi.Mux) *http.Server {

	hs := &http.Server{Addr: url, Handler: r}

	return hs
}
