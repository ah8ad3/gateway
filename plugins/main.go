package plugins

import (
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
)

// Plugin simple struct for plugins that define for service
type Plugin struct {
	Name   string
	Active bool
	Config map[string]interface{}
}

// Plugins all plugins that this gateway have and can set to proxies
var Plugins []Plugin

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

// FindPlugin to find and set plugin to proxy by costumer
// bool means error
func FindPlugin(name string, active bool, conf map[string]interface{}) (Plugin, bool) {

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
