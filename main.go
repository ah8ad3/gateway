package main

import (
	"fmt"
	"github.com/ah8ad3/gateway/routes"
	"net/http"
)

func main() {
	routes.LoadServices()
	routes.CheckServices()
	r := routes.Routes()
	fmt.Println("Server run at :3000")
	_ = http.ListenAndServe(":3000", r)
}
