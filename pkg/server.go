package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ah8ad3/gateway/pkg/integrate"

	"github.com/ah8ad3/gateway/pkg/db"
	exception "github.com/ah8ad3/gateway/pkg/err"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"github.com/ah8ad3/gateway/pkg/routes"
	"github.com/ah8ad3/gateway/plugins"
	"github.com/ah8ad3/gateway/plugins/auth"
	"github.com/go-chi/chi"
)

var serLocation, integrateLocation string

func init() {
	if os.Getenv("TEST") == "1" {
		serLocation = "./../services.json"
		integrateLocation = "./../integrates.json"
	} else {
		serLocation = "services.json"
		integrateLocation = "integrates.json"
	}
}

func settings() exception.Err {
	logger.OpenConnection()
	// Require for auth Database staff
	auth.OpenAuthCollection()

	// register all plugins to gateway
	plugins.RegisterPlugins()

	err := proxy.LoadServices(false, serLocation)
	if err.Message != "" {
		return err
	}
	proxy.CheckServices(false)

	err = integrate.LoadIntegration(integrateLocation)
	if err.Message != "" {
		return err
	}

	// check all service available every one hour
	go proxy.HealthCheck()
	return exception.Err{}
}

// RUN for run server
func RUN(ip string, port string, route string, stopSignal chan bool) {
	test := os.Getenv("TEST")

	str, _ := db.LoadSecretKey()
	db.SecretKey = str
	if db.SecretKey == "" {
		log.Fatal("Secret Key not Found, generate it by\n  \t gateway secret")
	}
	var handler *chi.Mux
	err := settings()
	if err.Critical {
		log.Fatal(err.Message)
	}

	if <-stopSignal {
		log.Fatal("StopCommand for server")
		// return

	}

	if route == "v1" {
		handler = routes.V1()
	} else {
		handler = routes.V2()
	}

	if port == "" {
		port = "3000"
	}
	if ip == "" {
		ip = "0.0.0.0"
	}
	listen := ip + ":" + port

	hs := setup(listen, handler)
	fmt.Printf("Listening on http://%s \n", hs.Addr)

	if test == "0" {
		_err := hs.ListenAndServe()
		if _err != http.ErrServerClosed {
			log.Fatal(_err.Error())
		}
	} else {
		fmt.Println("server ignored from listen and serve because of test mode if in production mode export, TEST=0")
	}
}

func setup(url string, handler *chi.Mux) *http.Server {

	hs := &http.Server{Addr: url, Handler: handler}

	return hs
}
