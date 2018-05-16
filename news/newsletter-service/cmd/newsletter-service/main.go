package main

import (
	"log"
	"newsletter-service/internal/db/mongodb"
	newsletterServer "newsletter-service/internal/server"
	"os"
	"os/signal"
)

//TODO: CLI mit urfave flags für dburl, dbname, collection name (constant values)
//TODO: architecture docu --> alla baustein sicht, gemeinsames dach aus arbeiten (alle), teams dann eigenes; Einfach nicht sonst was großes
//TODO: docu mit swagger, swaggerUI --> runterladen --> index files -> unter url= das .yaml angeben
//TODO: HTTPStatusCodes zurückgeben?
//TODO: Anleitung wie wird microservice gebastelt damit er zum schluss auch läuft
//TODO: Crosscompile mit go
//TODO: auf github migrieren mit repo
//TODO: POST returnen lassen

const (
	mongodbUrl            = "mongo:27017"
	mongodbName           = "newsletterService"
	mongodbCollectionName = "newsletterCollection"
	//TODO: constanten mit urfave als Flag machen (EnVar)
)

var mongoSession *mongodb.Session

func main() {
	// establish connection/session to mongodb
	mongoSession, err := mongodb.NewSession(mongodbUrl)
	if err != nil {
		log.Fatalln("Unable to connect to mongodb,. ERROR: '%s'", err)
	}

	// close session on newsletter-service shutdown
	defer mongoSession.Close()

	// service
	newsletterService := mongodb.NewNewsletterService(mongoSession.Copy(), mongodbName, mongodbCollectionName)

	// server
	server := newsletterServer.NewServer(newsletterService)
	server.Start()

	//gracefully shutdown
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
