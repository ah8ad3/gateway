package routes

type Url struct {
	Method string
	Path string
}

type Server struct {
	Server string
	Up bool
}

type Services struct {
	Name string
	Path string
	Server []Server
	Version int
	Description string
	Urls []Url
}
