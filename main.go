package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ah8ad3/gateway/pkg/auth"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/routes"
	"github.com/joho/godotenv"
)

func settings() {
	err := godotenv.Load()
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "auth", Time: time.Now()})
		log.Fatal("Error loading .env file")
	}
	logger.OpenConnection()
	// Require for auth Database staff
	auth.OpenAuthCollection()

	routes.LoadServices()
	routes.CheckServices(false)

	// check all service available every one hour
	go routes.HealthCheck()
}

func main() {
	settings()
	r := routes.V1()

	fmt.Println("Server run at :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "auth", Time: time.Now()})
		log.Fatal(err)
	}
}
