package ip

import "net/http"

var config = make(map[string]interface{}, 1)

// RegisterNewPlugin for add this plugin to plugin center
func RegisterNewPlugin() (string, bool, map[string]interface{}, func(config map[string]interface{}) func(handler http.Handler) http.Handler, int) {
	config = confTemplate(config)
	return "ipBlocker", true, config, Middleware, 1
}

func confTemplate(config map[string]interface{}) map[string]interface{}{
	config["nothing"] = nil

	return config
}
