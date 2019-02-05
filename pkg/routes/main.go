package routes

// V1Route created to work with V2 Routing
type V1Route interface {
	GetService() ([]byte, int)
	PostService() ([]byte, int)
	findService() string
}
