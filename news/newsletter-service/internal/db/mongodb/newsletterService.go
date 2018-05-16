package mongodb

import (
	"github.com/globalsign/mgo"
	"newsletter-service/internal/news"
	"github.com/globalsign/mgo/bson"
	"time"
)

type NewsletterService struct {
	collection *mgo.Collection
}

type NewsletterServiceHandler interface{
	CreateNewsletter(newsletter *news.Newsletter) error
	GetNewsletters() ([]*news.Newsletter, error)
	GetNewsletterById(mId string) (*news.Newsletter, error)
	DeleteNewsletterById(mId string) (error)
	GetUpcomingNewsletters() ([]*news.Newsletter, error)
}

func NewNewsletterService(mSession *Session, mDbName string, mCollectionName string) *NewsletterService{
	collection := mSession.GetCollection(mDbName, mCollectionName)
	return &NewsletterService{collection}
}

func(rNewsletterService *NewsletterService) CreateNewsletter(mNewsletter *news.Newsletter) (error) {
	return rNewsletterService.collection.Insert(&mNewsletter)
}

func(rNewsletterService *NewsletterService) GetNewsletters() ([]news.Newsletter, error){
	var newslettersSlice []news.Newsletter
	err := rNewsletterService.collection.Find(nil).All(&newslettersSlice)
	return newslettersSlice, err
}

func(rNewsletterService *NewsletterService) GetNewsletterById(mId string) (*news.Newsletter, error){
	newsletter := news.Newsletter{}
	err := rNewsletterService.collection.Find(bson.M{"_id":  bson.ObjectIdHex(mId)}).One(&newsletter)
	return &newsletter, err
}

func(rNewsletterService *NewsletterService) DeleteNewsletterById(mId string) error{
	err := rNewsletterService.collection.Remove(bson.M{"_id":  bson.ObjectIdHex(mId)})
	return err
}

func(rNewsletterService *NewsletterService) GetUpcomingEventNewsletters() ([]news.Newsletter, error){
	var newsletterSlice []news.Newsletter
	err := rNewsletterService.collection.Find(bson.M{"enddatum": bson.M{
		"$gt":time.Now(),
	}}).All(&newsletterSlice)

	return newsletterSlice, err
}
