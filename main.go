package main

var valid_methods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
var valid_protocol_versions = []string{"HTTP/1.1", "HTTP/1.0"}

func main() {

	// Example of how your application layer should ideally look:
	// router := NewRouter()
	// router.AddRoute("GET", "/", handleHome)
	// router.AddRoute("GET", "/about", handleAbout)
	// router.AddRoute("POST", "/users", handleCreateUser)
	// router.AddRoute("GET", "/users/:id", handleGetUser) // The hardest part

	// // The server takes the router and handles the TCP/parsing logic
	// server := NewServer(router)
	// server.Listen(":8080")

}

// func handleHome()
