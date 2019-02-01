package ip

import (
	"encoding/json"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/logger"
	"io/ioutil"
	"net/http"
	"time"
)

var blockList []BlockIpList
var apiIp []ApiIp

func init() {
	go updateBlockList()
}

func AddBlockList(ip string, path string, days int, ever bool)  {
	blockList = append(blockList, &BlockIpList{ip: ip, createdTime: time.Now(),
		expireTime:time.Now().Add(time.Hour * time.Duration(24 * days)), path: path, ever: ever, active: true})
}

func updateBlockList()  {
	time.Sleep(time.Duration(time.Minute * 5))
	for _, val := range blockList {
		if val.expireTime.Before(time.Now()) {
			val.active = false
		}
	}
}

func checkApiBlock()  {

}

func getApi(api string) *ApiIp {
	apiInfo := &ApiIp{}
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

// IpInfoMiddleware this must not use like this must implement
func IpInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiInfo := getApi(r.RemoteAddr)
		if apiInfo.status == "success" {
			apiIp = append(apiIp, apiInfo)
		}else {
			logger.SetUserLog(logger.UserLog{Time: time.Now(), Ip: r.RemoteAddr, RequestUrl: r.URL.Path,
				Log: logger.Log{Event: "critical", Description: "api ip not respond correct response at this"}})

		}
		limiter := getVisitor(r.RemoteAddr)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
