package pkg

import (
	"fmt"
	"github.com/ah8ad3/gateway/pkg/integrate"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"github.com/ah8ad3/gateway/pkg/routes"
	"github.com/ah8ad3/gateway/plugins/auth"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
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

	proxy.LoadServices()
	proxy.CheckServices(false)

	integrate.LoadIntegration()

	// check all service available every one hour
	go proxy.HealthCheck()
}

// RUN for run server
func RUN(ip string, port string) {
	settings()
	r := routes.V1()

	if port == "" {
		port = "3000"
	}
	if ip == "" {
		ip = "localhost"
	}

	listen := ip + ":" + port

	fmt.Println("Server run at ", listen)
	if err := http.ListenAndServe(listen, r); err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "auth", Time: time.Now()})
		log.Fatal(err)
	}
}