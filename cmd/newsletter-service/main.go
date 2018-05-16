package main

import (
	"log"
	"newsletter-service/internal/db/mongodb"
	newsletterServer "newsletter-service/internal/server"
	"os"
	"os/signal"

	"time"

	"github.com/urfave/cli"
)

//TODO: architecture docu --> alla baustein sicht, gemeinsames dach aus arbeiten (alle), teams dann eigenes; Einfach nicht sonst was großes
//TODO: docu mit swagger, swaggerUI --> runterladen --> index files -> unter url= das .yaml angeben
//TODO: Anleitung wie wird microservice gebastelt damit er zum schluss auch läuft
//TODO: POST returnen lassen / HTTPStatusCodes zurückgeben?

var mongoEnv *MongoEnv

type MongoEnv struct {
	session      *mongodb.Session
	host         string
	port         string
	dbName       string
	dbCollection string
}

func main() {
	mongoEnv = &MongoEnv{}

	/*
		Setup CLI
	*/
	app := cli.NewApp()
	app.Name = "newsletterService"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name: "Oliver Loose",
		},
		{
			Name: "Ricardo Zirk",
		},
	}
	// define flags/parameters
	app.Usage = "Start and manage the newsletter micro service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "db-host, dbh",
			Value:       "127.0.0.1",
			Usage:       "Set host of mongodb service",
			Destination: &mongoEnv.host,
		},
		cli.StringFlag{
			Name:        "db-port, dbp",
			Value:       ":27017",
			Usage:       "Set port of mongodb service",
			Destination: &mongoEnv.port,
		},
		cli.StringFlag{
			Name:        "db-name, dbn",
			Value:       "newsletterServiceDB",
			Usage:       "Set name of database to use",
			Destination: &mongoEnv.dbName,
		},
		cli.StringFlag{
			Name:        "db-collection, dbc",
			Value:       "newsletterCollection",
			Usage:       "Set name of collection within database to use",
			Destination: &mongoEnv.dbCollection,
		},
	}
	// define commands
	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Starts the newsletter micro service",
			Action: StartNewsletterServer,
		},
	}

	// run cli
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Error during start up: '%s'", err)
	}
}

func StartNewsletterServer(c *cli.Context) error {
	var err error

	// establish connection/session to mongodb
	log.Printf("Trying to connect to mongodb on '%s' with database '%s' and collection '%s'",
		mongoEnv.host+mongoEnv.port, mongoEnv.dbCollection, mongoEnv.dbCollection)

	mongoEnv.session, err = mongodb.NewSession(mongoEnv.host + mongoEnv.port)
	if err != nil {
		log.Fatalf("Unable to connect to mongodb: '%s'", err)
	} else {
		log.Print("Connected to database")
	}

	// close session on newsletter-service shutdown
	defer mongoEnv.session.Close()

	// setup newsletter service for database transactions
	newsletterService := mongodb.NewNewsletterService(mongoEnv.session.Copy(), mongoEnv.dbName,
		mongoEnv.dbCollection)

	// setup server
	server := newsletterServer.NewServer(newsletterService)
	server.Start()

	// gracefully shutdown
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt)
	signal.Notify(gracefulStop, os.Kill)
	go func() {
		<-gracefulStop

		if mongoEnv.session != nil {
			mongoEnv.session.Close()
		}

		os.Exit(0)
	}()

	return err
}
