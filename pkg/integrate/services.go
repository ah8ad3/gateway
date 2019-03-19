package integrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	exception "github.com/ah8ad3/gateway/pkg/err"
	"github.com/ah8ad3/gateway/pkg/logger"
)

const (
	// WarningColor yellow color
	WarningColor = "\033[1;33m%s\033[0m"
)

// Integrates to set all aggregations patterns
var Integrates []Integrate

// LoadIntegration method for load all aggregations from file
func LoadIntegration(location string) exception.Err  {
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return exception.Err{Message: err.Error(), Critical: true}.Log("system")
	}
	err = json.Unmarshal(data, &Integrates)
	if err != nil {
		return exception.Err{Message: err.Error(), Critical: true}.Log("system")

	}

	for _, val := range Integrates {
		fmt.Printf(WarningColor, fmt.Sprintf("Aggregate %s Loaded \n", val.Name))
	}
	return exception.Err{}
}

// GetIntegrateService function for implement GET at V1 Routing
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
