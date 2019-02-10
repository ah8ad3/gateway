package routes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/ah8ad3/gateway/pkg/proxy"

	"github.com/ah8ad3/gateway/pkg/logger"
)

const (
	// WarningColor yellow color
	WarningColor = "\033[1;33m%s\033[0m"

	// ErrorColor red color
	ErrorColor = "\033[1;31m%s\033[0m"
)

// GetService function for implement GET at V1 Routing
func GetService(server string, path string, query string) ([]byte, int) {
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

		return []byte(`{"error": "Service is Down!"}`), 404
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, 200
}

// PostService function for implement POST at V1 Routing
func PostService(server string, path string, query []byte) ([]byte, int) {
	url := "http://" + server + path
	req, _ := http.NewRequest("POST", url, bytes.NewReader(query))
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Description: "Service is down!", Event: "critical"},
			Time: time.Now(), RequestURL: url})
		return []byte(`{"error": "Service is Down!"}`), 404
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, 200
}

func findService(path string) string {
	rand.Seed(time.Now().Unix())
	path = "/" + path
	for _, val := range proxy.Services {
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
		Description: fmt.Sprintf("bad path check services %s", path)},
		Pkg: "auth", Time: time.Now()})
	//log.Fatal("bad path check services ", path)
	return ""
}
