package main

import (
	"fmt"
	"github.com/ah8ad3/gateway/logger"
	"github.com/ah8ad3/gateway/routes"
	"net/http"
)

func main() {
	logger.OpenConnection()

	routes.LoadServices()
	routes.CheckServices()
	r := routes.V1()
	fmt.Println("Server run at :3000")
	_ = http.ListenAndServe(":3000", r)
}
