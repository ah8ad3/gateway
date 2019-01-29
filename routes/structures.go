package routes

type Url struct {
	Method string
	Path string
}

type Services struct {
	Name string
	Path string
	Server string
	Version int
	Description string
	Urls []Url
}
