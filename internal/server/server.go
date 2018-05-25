package server

import (
	"log"
	"net/http"
	"newsletter-service/internal/db/mongodb"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/file"
	"github.com/go-ozzo/ozzo-routing/slash"
)

const (
	port = ":8080"
)

// Defines a Server type containing a router that defines the rest-api routing.
type Server struct {
	router *routing.Router
}

// Returns a server instance.
func NewServer(mNewsletterService *mongodb.NewsletterService) *Server {
	server := Server{router: routing.New()}
	// define base router configuration
	server.router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	// serve swagger-ui
	server.router.Get("/swagger/*", file.Server(file.PathMap{
		"/swagger/": "/swagger-ui/dist/",
	}))

	// add sup routes (/news/*)
	NewNewsRouter(mNewsletterService, &server)

	return &server
}

// Returns a new sub router for the referenced server instance.
func (rServer *Server) NewSubrouter(mPath string) *routing.RouteGroup {
	return rServer.router.Group(mPath)
}

// Starts a referenced server.
func (rServer *Server) Start() {
	log.Println("Listining on port :8080")
	// add api routes
	http.Handle("/", rServer.router)
	// start server
	http.ListenAndServe(port, nil)
}
