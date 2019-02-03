package ip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
)

var blockList []BlockIPList
var apiIP []APIIP

// AddBlockList to add some ip for blocking with expire time
func AddBlockList(ip string, path string, duration time.Duration, ever bool) {
	blockList = append(blockList, BlockIPList{ip: ip, createdTime: time.Now(),
		expireTime: time.Now().Add(duration), path: path, ever: ever, active: true})
}

// UpdateBlockList for update all block list for delete expired
func UpdateBlockList() {
	for {
		time.Sleep(time.Duration(time.Minute * 1))
		for listID, val := range blockList {
			if val.expireTime.Before(time.Now()) {
				blockList[listID].active = false
			}
		}
	}
}

func isAPIBlock(path string, ip string) bool {
	for _, val := range blockList {
		if val.active && val.path == path && val.ip == ip {
			return true
		}
	}
	return false
}

func getAPI(api string) *APIIP {
	apiInfo := &APIIP{}
	response, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", api))
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Time: time.Now(), Pkg: "plugins/ip",
			Log: logger.Log{Description: err.Error(), Event: "critical"}})
		return apiInfo
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(contents, apiInfo)
	return apiInfo
}

// InfoMiddleware this must not use like this must implement
func InfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		splitRoute := strings.Split(r.URL.Path, "/")
		// extract server path from url
		path := splitRoute[1]

		// This method disabled now but must implement with goroutine later
		//apiInfo := getAPI(r.RemoteAddr)
		apiInfo := &APIIP{status: "fail"}
		if apiInfo.status == "success" {
			apiIP = append(apiIP, *apiInfo)
			if isAPIBlock(path, r.RemoteAddr) {
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}
		} else {
			logger.SetUserLog(logger.UserLog{Time: time.Now(), IP: r.RemoteAddr, RequestURL: r.URL.Path,
				Log: logger.Log{Event: "critical", Description: "api ip not respond correct response at this"}})
			if isAPIBlock(path, r.RemoteAddr) {
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
