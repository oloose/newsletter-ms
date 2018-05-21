package server

import (
	"newsletter-service/internal/db/mongodb"
	"newsletter-service/internal/news"

	"log"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
)

type NewsRouter struct {
	newsletterService mongodb.NewsletterService
}

func NewNewsRouter(mNewsletterService *mongodb.NewsletterService, mRouter *routing.RouteGroup) *routing.RouteGroup {
	newsRouter := NewsRouter{*mNewsletterService}

	mRouter.Use(
		content.TypeNegotiator(content.JSON),
	)

	mRouter.Get("/newsletters", newsRouter.GetNewsletters, fault.Recovery(log.Printf))
	mRouter.Get(`/newsletter/<id>`, newsRouter.GetNewsletterById, fault.Recovery(log.Printf),
		fault.ErrorHandler(log.Printf)).Delete(newsRouter.DeleteNewsletterById, fault.Recovery(log.Printf))
	mRouter.Post("/newsletter", newsRouter.PostNewsletter, fault.Recovery(log.Printf)).Put(
		newsRouter.PutNewsletter)
	mRouter.Get("/newsletters/upcoming", newsRouter.GetUpcomingNewsletters, fault.Recovery(log.Printf))

	return mRouter
}

func (rNewsRouter *NewsRouter) GetNewsletters(mContext *routing.Context) error {
	// get newsletters from db
	newsletters, err := rNewsRouter.newsletterService.GetNewsletters()

	if err != nil {
		return err
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
		return err
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
	newNewsletter, err := newsletterParseObject.Parse()
	if err != nil {
		return routing.NewHTTPError(400, "Invalid input. (ERROR: "+err.Error()+")")
	}

	// store in db
	if err := rNewsRouter.newsletterService.CreateNewsletter(newNewsletter); err != nil {
		return err
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
		return err
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
		return err
	}

	return nil
}

func (rNewsRouter *NewsRouter) GetUpcomingNewsletters(mContext *routing.Context) error {
	// get upcoming (end date after now) newsletters from db
	newsletters, err := rNewsRouter.newsletterService.GetUpcomingNewsletters()

	if err != nil {
		return err
	}

	mContext.Write(newsletters)
	// check if list of newsletters is empty
	if len(newsletters) == 0 {
		return routing.NewHTTPError(204, "Result is empty. No newsletters available.")
	}

	return nil
}
