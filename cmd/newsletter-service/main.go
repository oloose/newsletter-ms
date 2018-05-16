package main

import (
	"log"
	"newsletter-service/internal/db/mongodb"
	newsletterServer "newsletter-service/internal/server"
	"os"
	"os/signal"
)

const (
	mongodbUrl            = "127.0.0.1:27017"
	mongodbName           = "newsletterService"
	mongodbCollectionName = "newsletterCollection"
	//TODO: constanten mit urfave als Flag machen
	//EnVar
)

var mongoSession *mongodb.Session

func main() {
	// establish connection/session to mongodb
	mongoSession, err := mongodb.NewSession(mongodbUrl)
	if err != nil {
		log.Fatalln("Unable to connect to mongodb")
	}

	// close session on newsletter-service shutdown
	defer mongoSession.Close()

	// service
	newsletterService := mongodb.NewNewsletterService(mongoSession.Copy(), mongodbName, mongodbCollectionName)

	// server
	server := newsletterServer.NewServer(newsletterService)
	server.Start()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt)
	signal.Notify(gracefulStop, os.Kill)

	go func() {
		<-gracefulStop

		if mongoSession != nil {
			mongoSession.Close()
		}

		os.Exit(0)
	}()
}
