package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/logger"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)
var Service []Services


func LoadServices()  {
	data, err := ioutil.ReadFile("services.json")
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "auth", Time: time.Now()})
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
		for _, server := range val.Server{
			if _, err := net.Dial("tcp", server.Server); err != nil{
				server.Up = false
				fmt.Printf(ErrorColor, fmt.Sprintf("Service %s not Up in server %s \n", val.Name, server.Server))
				logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "warning", Description: err.Error()},
					Pkg: "auth", Time: time.Now()})
				//log.Fatal(err)  // for production mode
			}else {
				server.Up = true
			}
		}
	}
}

func GetService(server string, path string, query string) []byte {
	url := "http://" + server + path
	if query != "" {
		url = url + "?" + query
	}
	req, _ :=http.NewRequest("GET", url, nil)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Description: "Service is down!", Event: "critical"},
			Time: time.Now(), RequestUrl: url})

		return []byte(`{"error": "Service is Down!"}`)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}


func PostService(server string, path string, query []byte) []byte {
	url := "http://" + server + path
	req, _ :=http.NewRequest("POST", url, bytes.NewReader(query))
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Description: "Service is down!", Event: "critical"},
			Time: time.Now(), RequestUrl: url})
		return []byte(`{"error": "Service is Down!"}`)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func findService(path string) string {
	rand.Seed(time.Now().Unix())
	path = "/" + path
	for _, val := range Service {
		if val.Path == path {
			for range val.Server{
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
