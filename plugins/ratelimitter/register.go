package ratelimitter

import "net/http"

var config = make(map[string]interface{}, 1)

// RegisterNewPlugin for add this plugin to plugin center
func RegisterNewPlugin() (string, bool, map[string]interface{}, func(config map[string]interface{}) func(handler http.Handler) http.Handler, int) {
	config = confTemplate(config)
	return "rateLimiter", true, config, Middleware, 2
}

func confTemplate(config map[string]interface{}) map[string]interface{}{
	config["block_time"] = 60
	config["rps"] = 5

	return config
}
