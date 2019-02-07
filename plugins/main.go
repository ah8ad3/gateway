package plugins

import (
	"github.com/ah8ad3/gateway/plugins/ip"
	"log"
	"net/http"
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/plugins/ratelimitter"
)

// Plugin simple struct for plugins that define for service
type Plugin struct {
	Name   string
	Active bool
	Config map[string]interface{}
	Middleware func(config map[string]interface{}) func(handler http.Handler) http.Handler `json:"-"`
	Priority int
}

// Plugins all plugins that this gateway have and can set to proxies
var Plugins []Plugin

// plugs unregistered plugins are here
var plugs []interface{}

// SetUPPlugin to add plugin fro creator to plugin list
// return error if true mean that error happened
func (p Plugin) SetUPPlugin() bool {
	if p.Name == "" {
		logger.SetSysLog(logger.SystemLog{Time: time.Now(), Pkg: "plugins", Log: logger.Log{Event: "log",
			Description: "plugin name is empty by plugin creator"}})
		return true
	}
	Plugins = append(Plugins, p)
	return false
}

// AddPluginProxy to find and set plugin to proxy by costumer
// bool means error
func AddPluginProxy(name string, active bool, conf map[string]interface{}) (Plugin, bool) {
	if name == "" {
		return Plugin{}, true
	}

	for _, val := range Plugins {
		if val.Name == name {
			if conf == nil {
				conf = val.Config
			}
			return Plugin{Name: val.Name, Active: active, Config: conf, Middleware: val.Middleware}, false
		}
	}

	return Plugin{}, true
}

// RegisterPlugins for add plugs to Plugins
func RegisterPlugins() {
	// add your'e plugin here
	plugs = append(plugs, ratelimitter.RegisterNewPlugin)
	plugs = append(plugs, ip.RegisterNewPlugin)

	for id := range plugs {
		name, active, config, middle, priority := plugs[id].(func()(string, bool, map[string]interface{}, func(config map[string]interface{}) func(handler http.Handler) http.Handler, int))()

		err := Plugin{Name: name, Active: active, Config: config, Middleware: middle, Priority: priority}.SetUPPlugin()

		if err{
			log.Fatal("name of plugin cant be empty")
		}
	}
}

// GetMiddleware use in SyncMiddleware for proxy after load of db
func GetMiddleware(pluginName string) func(config map[string]interface{}) func(handler http.Handler) http.Handler {
	for _, val := range Plugins {
		if val.Name == pluginName {
			return val.Middleware
		}
	}
	return nil
}
