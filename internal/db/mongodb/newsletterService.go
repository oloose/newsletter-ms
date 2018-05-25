// Package providing functions to establish a mongodb session and access the database.
package mongodb

import (
	"newsletter-service/internal/news"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Defines the negotiator for mongodb access which will be used as referencer type for functions that access/manipulate
// the database.
// Contains the accessed mongodb collection.
type NewsletterService struct {
	collection *mgo.Collection
}

// Interface defining service functions
type NewsletterServiceHandler interface {
	CreateNewsletter(newsletter *news.Newsletter) error
	UpdateNewsletter(newsletter *news.Newsletter) (*mgo.ChangeInfo, error)
	GetNewsletters() ([]*news.Newsletter, error)
	GetNewsletterById(mId string) (*news.Newsletter, error)
	DeleteNewsletterById(mId string) error
	GetUpcomingNewsletters() ([]*news.Newsletter, error)
}

// Returns a new NewsletterService instance with a initialized collection.
// Gets the collection from given Session (mSession) type based on the name of the mongodb database (mDbName)
// and collection (mCollectionName) to use.
func NewNewsletterService(mSession *Session, mDbName string, mCollectionName string) *NewsletterService {
	// get collection
	collection := mSession.GetCollection(mDbName, mCollectionName)
	return &NewsletterService{collection}
}

// Creates a new Newsletter entry in mongodb collection.
func (rNewsletterService *NewsletterService) CreateNewsletter(mNewsletter *news.Newsletter) error {
	return rNewsletterService.collection.Insert(&mNewsletter)
}

// Updates a Newsletter entry in mongodb based on the id of the given newsletter (mNewsletter).
func (rNewsletterService *NewsletterService) UpdateNewsletter(mNewsletter *news.Newsletter) (*mgo.ChangeInfo, error) {
	return rNewsletterService.collection.UpsertId(&mNewsletter.Id, &mNewsletter)
}

// Returns all newsletter in the mongodb collection.
func (rNewsletterService *NewsletterService) GetNewsletters() ([]news.Newsletter, error) {
	var newslettersSlice []news.Newsletter
	// find all newsletters and store in newsletterSlice
	err := rNewsletterService.collection.Find(nil).All(&newslettersSlice)
	return newslettersSlice, err
}

// Returns the newsletter with the given id (mId).
func (rNewsletterService *NewsletterService) GetNewsletterById(mId string) (*news.Newsletter, error) {
	var newsletter []news.Newsletter
	// try to find the newsletter based on mId and return it
	err := rNewsletterService.collection.FindId(bson.ObjectIdHex(mId)).All(&newsletter)
	return &newsletter[0], err
}

// Deletes a newsletter with the give id (mId).
func (rNewsletterService *NewsletterService) DeleteNewsletterById(mId string) error {
	return rNewsletterService.collection.Remove(bson.M{"_id": bson.ObjectIdHex(mId)})
}

// Returns all newsletter where the end date is in the future.
func (rNewsletterService *NewsletterService) GetUpcomingNewsletters() ([]news.Newsletter, error) {
	var newsletterSlice []news.Newsletter
	// find all newsletter with a enddatum date value in the future
	err := rNewsletterService.collection.Find(bson.M{"enddatum": bson.M{
		"$gt": time.Now(),
	}}).All(&newsletterSlice)

	return newsletterSlice, err
}
