package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
)

const (
	// WarningColor yellow color
	WarningColor = "\033[1;33m%s\033[0m"

	// ErrorColor red color
	ErrorColor = "\033[1;31m%s\033[0m"
)

// Services define all services that connect to gateway
var Services []Service

// LoadServices function for loading services from json file
func LoadServices() {
	data, err := ioutil.ReadFile("services.json")
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "routes", Time: time.Now()})
		log.Fatal(err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &Services)
	if err != nil {
		fmt.Println("services.json cant match to Structure read the docs or act like template")
		os.Exit(1)
	}

	for _, val := range Services {
		fmt.Printf(WarningColor, fmt.Sprintf("Service %s Loaded \n", val.Name))
	}
}

// CheckServices function for check if service is available or not
// params: then : mean that function call at start or later in health
func CheckServices(then bool) {
	for serviceID, val := range Services {
		for serverID, server := range val.Server {
			if _, err := net.Dial("tcp", server.Server); err != nil {
				Services[serviceID].Server[serverID].Up = false
				if then == false {
					fmt.Printf(ErrorColor, fmt.Sprintf("Service %s not Up in server %s \n", val.Name, server.Server))
				}
				logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "warning", Description: err.Error()},
					Pkg: "routes", Time: time.Now()})
				//log.Fatal(err)  // for production mode
			} else {
				fmt.Println("server ", server.Server, " is up")
				Services[serviceID].Server[serverID].Up = true
			}
		}
	}
}

// GetService function for implement GET at V1 Routing
func GetService(server string, path string, query string) []byte {
	url := "http://" + server + path
	if query != "" {
		url = url + "?" + query
	}
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Description: "Service is down!", Event: "critical"},
			Time: time.Now(), RequestURL: url})

		return []byte(`{"error": "Service is Down!"}`)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// PostService function for implement POST at V1 Routing
func PostService(server string, path string, query []byte) []byte {
	url := "http://" + server + path
	req, _ := http.NewRequest("POST", url, bytes.NewReader(query))
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Description: "Service is down!", Event: "critical"},
			Time: time.Now(), RequestURL: url})
		return []byte(`{"error": "Service is Down!"}`)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func findService(path string) string {
	rand.Seed(time.Now().Unix())
	path = "/" + path
	for _, val := range Services {
		if val.Path == path {
			for range val.Server {
				ser := val.Server[rand.Intn(len(val.Server))]
				if ser.Up {
					return ser.Server
				}
			}
			return val.Server[rand.Intn(len(val.Server))].Server
		}
	}
	logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical",
		Description: fmt.Sprintf("bad path check services ", path)},
		Pkg: "auth", Time: time.Now()})
	//log.Fatal("bad path check services ", path)
	return ""
}

// HealthCheck function for check all services per hour in goroutine
func HealthCheck() {
	for {
		time.Sleep(time.Duration(time.Hour * 1))
		CheckServices(true)
	}
}

func getServices() []byte {
	var jData []byte
	jData, _ = json.Marshal(Services)

	return jData
}
