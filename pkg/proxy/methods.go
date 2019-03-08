package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
func LoadServices(jsonData bool, serLocation string) error {
	if jsonData {
		data, err := ioutil.ReadFile(serLocation)

		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &Services)
		if err != nil {
			//fmt.Println("services.json cant match to Structure read the docs or act like template")
			//os.Exit(1)
			return err
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

	return nil
}

// CheckServices function for check if service is available or not
// params: then : mean that function call at start or later in health
func CheckServices(then bool) {
	for serviceID, val := range Services {
		upHostsCount := 0
		var upHosts []string
		for serverID, server := range val.Server {

			timeout := time.Duration(3 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			_, err := client.Get(server.Server + "/healthz")
			if err != nil {
				Services[serviceID].Server[serverID].Up = false
				if then == false {
					fmt.Printf(ErrorColor, fmt.Sprintf("Service %s not Up in server %s \n", val.Name, server.Server))
				}
				logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "warning", Description: err.Error()},
					Pkg: "routes", Time: time.Now()})
				//log.Fatal(err)  // for production mode
			} else {
				upHostsCount++
				upHosts = append(upHosts, server.Server)
				Services[serviceID].Server[serverID].Up = true
			}
		}
		Services[serviceID].UPHostsCoutn = upHostsCount
		Services[serviceID].UPHosts = upHosts
	}
}

// HealthCheck function for check all services per hour in goroutine
func HealthCheck() {
	for {
		time.Sleep(time.Duration(time.Minute * 1))
		CheckServices(true)
	}
}

func updateConfigPlugin(pluginName string, serviceID int, pluginID int, config map[string]interface{}) {
	// search if one of the configs not entered by user be default value
	for _, plug := range plugins.Plugins {
		if pluginName == plug.Name {
			for key, value := range plug.Config {
				if config[key] == nil {
					config[key] = value
				}
			}
			Services[serviceID].Plugins[pluginID].Config = config
			// save services after change automatically
			SaveServices()
		}
	}
}

// AddPlugin api for add plugin to proxy
// work with service name and version and plugin name
func AddPlugin(serviceName string, version int, pluginName string, config map[string]interface{}) (string, bool) {
	for id, val := range Services {
		if val.Name == serviceName && val.Version == version {
			for idx, plug := range val.Plugins {
				if plug.Name == pluginName {
					if config != nil {
						updateConfigPlugin(pluginName, id, idx, config)
						return "config updated", false
					}
					return "plugin exist", true
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
func RemovePlugin(serviceName string, version int, pluginName string) (string, bool) {
	for id, val := range Services {
		if serviceName == val.Name && version == val.Version {
			for idx, plug := range val.Plugins {
				if plug.Name == pluginName {
					Services[id].Plugins = append(val.Plugins[:idx], val.Plugins[idx+1:]...)

					SaveServices()
					return "Deleted", false
				}
			}
			return "Plugin not found", true
		}
	}
	return "Proxy not found", true
}

// SaveServices to save services
// you just call it and this function save all services are in proxy
func SaveServices() {
	JData, _ := json.Marshal(Services)
	db.InsertProxy(JData)
}
