package pkg

import (
	"context"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/plugins"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ah8ad3/gateway/pkg/integrate"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"github.com/ah8ad3/gateway/pkg/routes"
	"github.com/ah8ad3/gateway/plugins/auth"
	"github.com/joho/godotenv"
)

func settings() {
	err := godotenv.Load()
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "main", Time: time.Now()})
		log.Fatal("Error loading .env file")
	}
	logger.OpenConnection()
	// Require for auth Database staff
	auth.OpenAuthCollection()

	// register all plugins to gateway
	plugins.RegisterPlugins()

	proxy.LoadServices(false)
	proxy.CheckServices(false)

	integrate.LoadIntegration()

	// check all service available every one hour
	go proxy.HealthCheck()
}

// RUN for run server
func RUN(ip string, port string) {
	str, _ := db.LoadSecretKey()
	db.SecretKey = str
	if db.SecretKey == "" {
		log.Fatal("Secret Key not Found, generate it by\n  \t gateway secret")
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	settings()
	r := routes.V1()

	if port == "" {
		port = "3000"
	}
	if ip == "" {
		ip = "localhost"
	}

	listen := ip + ":" + port

	hs := setup(listen, r)

	go func() {
		fmt.Println(fmt.Sprintf("Listening on http://%s\n", hs.Addr))

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-stop

	_ = hs.Shutdown(context.Background())
}

func setup(url string, r *chi.Mux) *http.Server {

	hs := &http.Server{Addr: url, Handler: r}

	return hs
}
