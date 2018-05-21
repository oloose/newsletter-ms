package news

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

type Newsletter struct {
	Id                   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Beschreibung         string        `json:"beschreibung"`
	BeschreibungEnglisch string        `json:"beschreibungEnglisch"`
	Enddatum             time.Time     `json:"enddatum"`
	Person               string        `json:"person"`
	Startdatum           time.Time     `json:"startdatum"`
	Titel                string        `json:"titel"`
	TitelEnglisch        string        `json:"titelEnglisch"`
	Verdatum             time.Time     `json:"verdatum"`
}

type NewsletterParseObject struct {
	Id                   string
	Beschreibung         string
	BeschreibungEnglisch string
	Enddatum             string
	Person               string
	Startdatum           string
	Titel                string
	TitelEnglisch        string
	Verdatum             string
}

func (rNewsletterParseObject *NewsletterParseObject) Parse() (*Newsletter, error) {
	// copy possible values in newsletter object
	newsletter := Newsletter{
		Beschreibung:         rNewsletterParseObject.Beschreibung,
		BeschreibungEnglisch: rNewsletterParseObject.BeschreibungEnglisch,
		// Enddatum: nil,
		Person: rNewsletterParseObject.Person,
		// Startdatum: nil,
		Titel:         rNewsletterParseObject.Titel,
		TitelEnglisch: rNewsletterParseObject.TitelEnglisch,
		// Verdatum: nil,
	}

	// parse date from string to time.Time
	location, err := time.LoadLocation("Europe/Berlin")
	newsletter.Enddatum, err = time.ParseInLocation("02.01.06; 15:04", rNewsletterParseObject.Enddatum,
		location)
	newsletter.Startdatum, err = time.ParseInLocation("02.01.06; 15:04", rNewsletterParseObject.Startdatum,
		location)
	newsletter.Verdatum, err = time.ParseInLocation("02.01.06; 15:04", rNewsletterParseObject.Verdatum,
		location)

	// check for existing id
	if rNewsletterParseObject.Id != "" {
		newsletter.Id = bson.ObjectIdHex(rNewsletterParseObject.Id)
	}

	if err != nil {
		return &newsletter, errors.New("Wrong date format! Error while parsing date. " +
			"Use date format: '02.01.06; 15:04' (" + err.Error() + ")")
	}

	return &newsletter, nil
}
