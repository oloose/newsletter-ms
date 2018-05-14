package server

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"log"
	"github.com/go-ozzo/ozzo-routing/slash"
	"net/http"
	"github.com/go-ozzo/ozzo-routing/fault"
	"newsletter-service/internal/db/mongodb"
)

const(
	port = ":8080"
)

type Server struct {
	router *routing.Router
}

func NewServer(mNewsletterService *mongodb.NewsletterService) *Server {
	server := Server{router: routing.New()}
	server.router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	// add sup routes
	NewNewsRouter(mNewsletterService, server.NewSubrouter("/news"))

	return &server
}

func(rServer *Server) NewSubrouter(mPath string) *routing.RouteGroup{
	return rServer.router.Group(mPath)
}

func(rServer *Server) Start(){
	log.Println("Listining on port :8080")
	http.Handle("/", rServer.router)
	http.ListenAndServe(port, nil)
}

