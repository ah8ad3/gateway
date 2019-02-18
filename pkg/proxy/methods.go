package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/plugins"

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
func LoadServices(jsonData bool, serLocation string) {
	if jsonData {
		data, err := ioutil.ReadFile(serLocation)

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		err = json.Unmarshal(data, &Services)
		if err != nil {
			fmt.Println("services.json cant match to Structure read the docs or act like template")
			os.Exit(1)
		}

		// save data to db automatically after load
		SaveServices()
	} else {
		// this is how get info from db
		_ = json.Unmarshal(db.GetProxies(), &Services)
	}

	for _, val := range Services {
		SyncPlugins(val.Name)
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

// HealthCheck function for check all services per hour in goroutine
func HealthCheck() {
	for {
		time.Sleep(time.Duration(time.Hour * 1))
		CheckServices(true)
	}
}

// AddPlugin api for add plugin to proxy
func AddPlugin(serviceName string, pluginName string, config map[string]interface{}) (string, bool) {
	for id, val := range Services {
		if val.Name == serviceName {
			for idx, plug := range val.Plugins {
				if plug.Name == pluginName {
					if config != nil {
						Services[id].Plugins[idx].Config = config
						return "ok", true
					}
					return "plugin exist", false
				}
			}

			plugin, err := plugins.AddPluginProxy(pluginName, true, config)
			if err {
				logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "log", Description: "plugin not found or name is empty"},
					Time: time.Now(), Pkg: "proxy"})
				return "plugin not found or name is empty", true
			}
			Services[id].Plugins = append(Services[id].Plugins, plugin)

			// save services after change automatically
			SaveServices()

			return "ok", false
		}
	}
	logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "log", Description: "proxy not found"},
		Time: time.Now(), Pkg: "proxy"})
	return "proxy not found", true
}

// SyncPlugins after load from db, sync functions to struct cause method address after app execution not meaning
// anything
func SyncPlugins(proxyName string) {
	for id, val := range Services {
		if val.Name == proxyName {
			for idx, plug := range val.Plugins {
				mid := plugins.GetMiddleware(plug.Name)
				if mid != nil {
					Services[id].Plugins[idx].Middleware = mid
				}
			}
		}
	}
}

// RemovePlugin from proxy
func RemovePlugin(serviceName string, pluginName string) (string, bool) {
	return "", false
}

// SaveServices to save services
// you just call it and this function save all services are in proxy
func SaveServices() {
	JData, _ := json.Marshal(Services)
	db.InsertProxy(JData)
}
