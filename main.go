package main

import (
	"fmt"
	"github.com/ah8ad3/gateway/pkg/auth"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func settings() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logger.OpenConnection()
	// Require for auth Database staff
	auth.OpenAuthCollection()

	routes.LoadServices()
	routes.CheckServices()

}

func main() {
	settings()
	r := routes.V1()
	fmt.Println("Server run at :3000")
	if err := http.ListenAndServe(":3000", r); err != nil{
		log.Fatal(err)
	}
}
