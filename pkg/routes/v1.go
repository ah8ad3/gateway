package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/integrate"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
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

func getProxyHttp(writer http.ResponseWriter, request *http.Request) {

	// remove path form url and send to service and serve answer
	splitRoute := strings.Split(request.URL.Path, "/")
	route := strings.Join(splitRoute[2:], "/")
	if route == "" {
		route = "/"
	} else {
		route = "/" + route
	}

	logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log"}, RequestURL: request.URL.Path,
		IP: request.RemoteAddr, Time: time.Now()})

	writer.Header().Set("Content-Type", "application/json")

	server := findService(splitRoute[1])
	body, code := GetService(server, route, request.URL.RawQuery)
	writer.WriteHeader(code)
	_, _ = writer.Write(body)
}

func postProxyHttp(writer http.ResponseWriter, request *http.Request) {
	// remove path form url and send to service and serve answer
	splitRoute := strings.Split(request.URL.Path, "/")
	route := strings.Join(splitRoute[2:], "/")
	if route == "" {
		route = "/"
	} else {
		route = "/" + route
	}
	_ = request.ParseForm()

	m := make(map[string]interface{})
	for key, value := range request.Form {
		m[key] = strings.Join(value, "")
	}

	data, _ := json.Marshal(m)

	logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log"}, RequestURL: request.URL.Path,
		IP: request.RemoteAddr, Time: time.Now()})

	writer.Header().Set("Content-Type", "application/json")
	server := findService(splitRoute[1])
	body, code := PostService(server, route, data)
	writer.WriteHeader(code)
	_, _ = writer.Write(body)
}

func putProxyHttp(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("hello"))
}

func deleteProxyHttp(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("hello"))
}

func integrateProxyHttp(w http.ResponseWriter, r *http.Request) {
	var result []map[string]interface{}
	_ = result
	for _, val := range integrate.Integrates {
		if val.Path == r.URL.Path {
			for _, service := range val.Join {
				var ser []map[string]interface{}
				_ = ser
				url := r.Host + service
				res, err := integrate.GetIntegrateService(url, r.URL.RawQuery)

				// check if service offline create error cause fixed aggregation
				if err && val.Fixed {
					logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log",
						Description: "One of the services was offline in aggregation"}, RequestURL: r.URL.Path,
						IP: r.RemoteAddr, Time: time.Now()})
					_, _ = w.Write([]byte(`{"error": "Aggregation failed one of the services are offline, log stored"}`))

					return
				}

				_ = json.Unmarshal(res, &ser)
				for _, item := range ser {
					result = append(result, item)
				}
			}
		}
	}
	// check if all the sevices are offline
	if result == nil {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Event: "log",
			Description: "all of the services was offline in aggregation"}, RequestURL: r.URL.Path,
			IP: r.RemoteAddr, Time: time.Now()})
		_, _ = w.Write([]byte(`{"error": "Aggregation failed all of the services are offline, log stored"}`))

		return
	}
	jData, _ := json.Marshal(result)
	_, _ = w.Write(jData)
	return
}
