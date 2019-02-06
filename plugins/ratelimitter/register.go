package ratelimitter

var config map[string]interface{}

// RegisterNewPlgin for add this plugin to plugin center
func RegisterNewPlgin() (string, bool, map[string]interface{}) {
	confTemplate()
	return "rateLimitter", true, config
}

func confTemplate() {
	config["block_time"] = 60
	config["rps"] = 5
}
