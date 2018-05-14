package main

import (
	"newsletter-service/internal/db/mongodb"
	"log"
	newsletterServer "newsletter-service/internal/server"
)

const(
	mongodbUrl = "127.0.0.1:27017"
	mongodbName = "newsletterService"
	mongodbCollectionName = "newsletterCollection"
)

func main() {
	// establish connection/session to mongodb
	mongoSession, err := mongodb.NewSession(mongodbUrl)
	if err != nil{
		log.Fatalln("Unable to connect to mongodb")
	}

	// close session on app shutdown
	defer mongoSession.Close()

	// service
	newsletterService := mongodb.NewNewsletterService(mongoSession.Copy(), mongodbName, mongodbCollectionName)

	// server
	server := newsletterServer.NewServer(newsletterService)
	server.Start()
}
