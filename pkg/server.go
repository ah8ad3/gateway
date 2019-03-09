package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ah8ad3/gateway/pkg/db"
	exception "github.com/ah8ad3/gateway/pkg/err"
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

func settings() exception.Err{
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

	//integrate.LoadIntegration()

	// check all service available every one hour
	go proxy.HealthCheck()
	return exception.Err{}
}

// RUN for run server
func RUN(ip string, port string, route string, test int) {
	str, _ := db.LoadSecretKey()
	db.SecretKey = str
	if db.SecretKey == "" {
		log.Fatal("Secret Key not Found, generate it by\n  \t gateway secret")
	}
	var r *chi.Mux
	err := settings()
	if err.Critical{
		log.Fatal(err.Message)
	}

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

	if test == 0 {
		_err := hs.ListenAndServe()
		if _err != http.ErrServerClosed {
			log.Fatal(_err.Error())
		}
	}else {
		fmt.Println("server ignored from listen and serve because of test mode if in production mode export, TEST=0")
	}
}

func setup(url string, r *chi.Mux) *http.Server {

	hs := &http.Server{Addr: url, Handler: r}

	return hs
}
