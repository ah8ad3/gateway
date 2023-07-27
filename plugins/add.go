package plugins

import "net/http"

type PluginAdder interface {
	// name, is_active, middleware
	AddPlugin(string, bool, func(handler http.Handler) http.Handler)
}
