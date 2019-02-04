package integrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

// Integrates to set all aggregations patterns
var Integrates []Integrate

// LoadIntegration method for load all aggregations from file
func LoadIntegration() {
	data, err := ioutil.ReadFile("integrates.json")
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "integrate", Time: time.Now()})
		log.Fatal(err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &Integrates)
	if err != nil {
		fmt.Println("integrates.json cant match to Structure read the docs or act like template")
		os.Exit(1)
	}

	for _, val := range Integrates {
		fmt.Printf(WarningColor, fmt.Sprintf("Aggregate %s Loaded \n", val.Name))
	}
}


// GetService function for implement GET at V1 Routing
func GetIntegrateService(url string, query string) ([]byte, bool) {
	url = "http://" + url
	if query != "" {
		url = url + "?" + query
	}

	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}

	res, _ := client.Do(req)

	if res.StatusCode != 200 {
		logger.SetUserLog(logger.UserLog{Log: logger.Log{Description: "Service is down!", Event: "critical"},
			Time: time.Now(), RequestURL: url})

		return []byte(`{"error": "Service is Down!"}`), true
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, false
}

