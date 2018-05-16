package server

import (
	"newsletter-service/internal/db/mongodb"
	"newsletter-service/internal/news"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
)

type NewsRouter struct {
	newsletterService mongodb.NewsletterService
}

func NewNewsRouter(mNewsletterService *mongodb.NewsletterService, mRouter *routing.RouteGroup) *routing.RouteGroup {
	newsRouter := NewsRouter{*mNewsletterService}

	mRouter.Use(
		content.TypeNegotiator(content.JSON),
	)

	mRouter.Get("/newsletters", newsRouter.GetNewsletters)
	mRouter.Get(`/newsletter/<id>`, newsRouter.GetNewsletterById).Delete(newsRouter.DeleteNewsletterById)
	mRouter.Post("/newsletter", newsRouter.PostNewsletter)
	mRouter.Get("/newsletters/upcoming", newsRouter.GetUpcomingNewsletters)

	return mRouter
}

func (rNewsRouter *NewsRouter) GetNewsletters(mContext *routing.Context) error {
	// get newsletters from db
	newsletters, err := rNewsRouter.newsletterService.GetNewsletters()

	if err != nil {
		return err
	}

	mContext.Write(newsletters)

	return nil
}

func (rNewsRouter *NewsRouter) GetNewsletterById(mContext *routing.Context) error {
	id := mContext.Param("id")
	// get from db
	newsletter, err := rNewsRouter.newsletterService.GetNewsletterById(id)
	if err != nil {
		return err
	}

	mContext.Write(newsletter)

	return nil
}

func (rNewsRouter *NewsRouter) PostNewsletter(mContext *routing.Context) error {
	var newsletterParseObject news.NewsletterParseObject
	err := mContext.Read(&newsletterParseObject)
	if err != nil {
		return err
	}

	// parse to newsletter
	newNewsletter, err := newsletterParseObject.Parse()
	if err != nil {
		return err
	}

	// store in db
	if err := rNewsRouter.newsletterService.CreateNewsletter(newNewsletter); err != nil {
		return err
	}

	return nil
}

func (rNewsRouter *NewsRouter) DeleteNewsletterById(mContext *routing.Context) error {
	id := mContext.Param("id")
	// delete from db
	err := rNewsRouter.newsletterService.DeleteNewsletterById(id)
	if err != nil {
		return err
	}

	return nil
}

func (rNewsRouter *NewsRouter) GetUpcomingNewsletters(mContext *routing.Context) error {
	// get upcoming (end date after now) newsletters from db
	newsletters, err := rNewsRouter.newsletterService.GetUpcomingEventNewsletters()

	if err != nil {
		return err
	}

	mContext.Write(newsletters)

	return nil
}
