package server

import (
	"newsletter-service/internal/db/mongodb"
	"newsletter-service/internal/news"

	"log"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
)

// Defines a NewsRouter type containing a NewsletterService to handle mongodb access.
type NewsRouter struct {
	newsletterService mongodb.NewsletterService
}

// Returns a NewsRouter as child route group of the base server (mServer) Router, that uses the given NewsletterService
// to handle mongodb access. Defines the REST-API routes that can be accessed.
func NewNewsRouter(mNewsletterService *mongodb.NewsletterService, mServer *Server) *routing.RouteGroup {
	newsRouter := NewsRouter{*mNewsletterService}
	// new sup router with /news/*
	newsSupRoute := mServer.NewSubrouter("/news")

	// routes only accept and return json data
	newsSupRoute.Use(
		content.TypeNegotiator(content.JSON),
	)

	// define REST routes
	newsSupRoute.Get("/newsletters", newsRouter.GetNewsletters, fault.Recovery(log.Printf))
	newsSupRoute.Get(`/newsletter/<id>`, newsRouter.GetNewsletterById, fault.Recovery(log.Printf)).Delete(
		newsRouter.DeleteNewsletterById, fault.Recovery(log.Printf))
	newsSupRoute.Post("/newsletter", newsRouter.PostNewsletter, fault.Recovery(log.Printf)).Put(
		newsRouter.PutNewsletter)
	newsSupRoute.Get("/newsletters/upcoming", newsRouter.GetUpcomingNewsletters, fault.Recovery(log.Printf))

	return newsSupRoute
}

func (rNewsRouter *NewsRouter) GetNewsletters(mContext *routing.Context) error {
	// get newsletters from db
	newsletters, err := rNewsRouter.newsletterService.GetNewsletters()

	if err != nil {
		return routing.NewHTTPError(500, err.Error())
	}

	mContext.Write(&newsletters)

	// check if list of newsletters is empty
	if len(newsletters) == 0 {
		return routing.NewHTTPError(204, "Result is empty. No newsletters available.")
	}
	return nil
}

func (rNewsRouter *NewsRouter) GetNewsletterById(mContext *routing.Context) error {
	id := mContext.Param("id")
	// get from db
	newsletter, err := rNewsRouter.newsletterService.GetNewsletterById(id)
	if err != nil {
		if err.Error() == "not found" {
			return routing.NewHTTPError(404, "Newsletter not found.(ERROR: "+err.Error()+")")
		}
		return routing.NewHTTPError(500, err.Error())
	}

	mContext.Write(newsletter)

	return nil
}

func (rNewsRouter *NewsRouter) PostNewsletter(mContext *routing.Context) error {
	var newsletterParseObject news.NewsletterParseObject
	err := mContext.Read(&newsletterParseObject)
	if err != nil {
		return routing.NewHTTPError(500, "(ERROR: "+err.Error()+")")
	}

	// parse to newsletter
	newsletterParseObject.Id = "" // post should create a new entry, so set id nil
	newNewsletter, err := newsletterParseObject.Parse()
	if err != nil {
		return routing.NewHTTPError(400, "Invalid input. (ERROR: "+err.Error()+")")
	}

	// store in db
	if err := rNewsRouter.newsletterService.CreateNewsletter(newNewsletter); err != nil {
		return routing.NewHTTPError(500, err.Error())
	}

	return nil
}

func (rNewsRouter *NewsRouter) PutNewsletter(mContext *routing.Context) error {
	var newsletterParseObject news.NewsletterParseObject
	err := mContext.Read(&newsletterParseObject)
	if err != nil {
		return routing.NewHTTPError(500, "(ERROR: "+err.Error()+")")
	}

	// parse to newsletter
	newNewsletter, err := newsletterParseObject.Parse()
	if err != nil {
		return routing.NewHTTPError(400, "Invalid input. (ERROR: "+err.Error()+")")
	}

	// store in db
	if _, err := rNewsRouter.newsletterService.UpdateNewsletter(newNewsletter); err != nil {
		if err.Error() == "not found" {
			return routing.NewHTTPError(404, "Newsletter not found.(ERROR: "+err.Error()+")")
		}
		return routing.NewHTTPError(500, err.Error())
	}

	return nil
}

func (rNewsRouter *NewsRouter) DeleteNewsletterById(mContext *routing.Context) error {
	id := mContext.Param("id")
	// delete from db
	err := rNewsRouter.newsletterService.DeleteNewsletterById(id)
	if err != nil {
		if err.Error() == "not found" {
			return routing.NewHTTPError(404, "Newsletter not found.(ERROR: "+err.Error()+")")
		}
		return routing.NewHTTPError(500, err.Error())
	}

	return nil
}

func (rNewsRouter *NewsRouter) GetUpcomingNewsletters(mContext *routing.Context) error {
	// get upcoming (end date after now) newsletters from db
	newsletters, err := rNewsRouter.newsletterService.GetUpcomingNewsletters()

	if err != nil {
		return routing.NewHTTPError(500, err.Error())
	}

	mContext.Write(newsletters)
	// check if list of newsletters is empty
	if len(newsletters) == 0 {
		return routing.NewHTTPError(204, "Result is empty. No newsletters available.")
	}

	return nil
}
