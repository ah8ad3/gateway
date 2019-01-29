package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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
		log.Fatal(err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &Service)
	if err != nil {
		fmt.Println("services.json cant match to Structure read the docs or act like template")
		os.Exit(1)
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

func GetService(server string, path string, method string, query string) []byte {
	url := "http://" + server + path
	if query != "" {
		url = url + "?" + query
	}
	req, _ :=http.NewRequest(method, url, nil)
	client := &http.Client{}

	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
