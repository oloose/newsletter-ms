// Package defining a Newsletter type and related functions.
package news

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// Defines a newsletter.
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

// Used as a placeholder while getting a Newsletter entry from the database, this object will latter parsed to a real
// Newsletter type object. This is necessary because all values in the returned newsletter entry from the database are
// strings (json) but Newsletter type contains time.Time attributes, so a direct conversion might throw an error.
// TODO: There might be a better approach to handle conversion we do not know yet.
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

// Parses all values from a referenced NewsletterParseObject into a newsletter type object and returns the newsletter.
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

	if err != nil {
		return &newsletter, errors.New("Wrong date format! Error while parsing date. " +
			"Use date format: '02.01.06; 15:04' (" + err.Error() + ")")
	}

	// check for existing id
	if rNewsletterParseObject.Id != "" {
		newsletter.Id = bson.ObjectIdHex(rNewsletterParseObject.Id)
	}

	return &newsletter, nil
}
