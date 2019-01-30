package main

import (
	"fmt"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/routes"
	"log"
	"net/http"
)

func settings() {
	logger.OpenConnection()

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
