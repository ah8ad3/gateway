package plugins

import (
	"log"
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/plugins/ratelimitter"
)

// Plugin simple struct for plugins that define for service
type Plugin struct {
	Name   string
	Active bool
	Config map[string]interface{}
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
			return Plugin{Name: val.Name, Active: active, Config: conf}, false
		}
	}

	return Plugin{}, true
}

// RegisterPlugins for add plugs to Plugins
func RegisterPlugins() {
	// add your'e plugin here
	plugs = append(plugs, ratelimitter.RegisterNewPlugin)

	for id := range plugs {
		name, active, config := plugs[id].(func()(string, bool, map[string]interface{}))()

		err := Plugin{Name: name, Active: active, Config: config}.SetUPPlugin()

		if err{
			log.Fatal("name of plugin cant be empty")
		}
	}
}