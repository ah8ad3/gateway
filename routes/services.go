package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

const (
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)
var Service []Services

func LoadServices()  {
	data, err := ioutil.ReadFile("services.json")
	if err != nil {
		fmt.Println("services.json file not found")
		os.Exit(2)
	}
	err = json.Unmarshal(data, &Service)
	if err != nil {
		fmt.Println("services.json cant match to Structure read the docs or act like template")
		os.Exit(2)
	}

	for _, val := range Service {
		fmt.Printf(WarningColor, fmt.Sprintf("Service %s Loaded \n", val.Name))
	}
}

func CheckServices() {
	for _, val := range Service {
		if _, err := net.Dial("tcp", val.Server); err != nil{
			fmt.Printf(ErrorColor, fmt.Sprintf("Service %s not Up \n", val.Name))
			//log.Fatal(err)  // for production mode
		}
	}
}
