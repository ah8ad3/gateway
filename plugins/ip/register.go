package ip

// RegisterNewPlugin for add this plugin to plugin center
func RegisterNewPlugin() (string, bool, map[string]interface{}) {
	config := make(map[string]interface{}, 1)
	confTemplate(config)
	return "rateLimiter", true, config
}

func confTemplate(config map[string]interface{}) map[string]interface{}{
	config["block_time"] = 60
	config["rps"] = 5

	return config
}
