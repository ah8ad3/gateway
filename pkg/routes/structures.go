package routes

// URL define routes of services
type URL struct {
	Method string
	Path   string
}

// Server define services active servers
type Server struct {
	Server string
	Up     bool
}

// Services define structure of service loaded to this gateway
type Services struct {
	Name        string
	Path        string
	Server      []Server
	Version     int
	Description string
	Urls        []URL
}
