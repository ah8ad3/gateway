package v1

// RequestForward model of request forwarding to do some customized aggregate by user or something
// Define some role to do some aggregate or filter in response
// Example:
// user GET /books
// Go to -> /log -> /api/book -> /filter and take the response back
type RequestForward struct {
	// Name describe name of forwarder like an identify
	Name string

	// BaseURL of start node
	BaseURL string

	// From the start node of forward
	From string
}

// func ()
