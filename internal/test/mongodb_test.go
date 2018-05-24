// Package containing mongodb related functions for testing (OUTDATED). Not everything possible is tested because this
// is just an example of how tests in golang could be done.
package mongodb

import (
	"log"
	"newsletter-service/internal/db/mongodb"
	"newsletter-service/internal/news"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
)

var TestSession mongodb.Session
var TestDummyNewsletters []news.Newsletter

const (
	mongodbUrl               = "localhost:27017"
	dbName                   = "test_newsletterdb"
	newsletterCollectionName = "test_NewsletterCollection"

	testBeschreibung         = "testBeschreibung"
	testBeschreibungEnglisch = "testBeschreibungEnglisch"
	testPerson               = "testPerson"
	testTitel                = "testTitel"
	testTitelEnglisch        = "testTitelEnglisch"
)

func Setup() {
	// setup dummy newsletters
	TestDummyNewsletters = append(TestDummyNewsletters, news.Newsletter{
		Id:                   bson.NewObjectId(),
		Beschreibung:         testBeschreibung,
		BeschreibungEnglisch: testBeschreibungEnglisch,
		Enddatum:             time.Now().AddDate(0, 0, 5),
		Person:               testPerson,
		Startdatum:           time.Now(),
		Titel:                testTitel,
		TitelEnglisch:        testTitelEnglisch,
		Verdatum:             time.Now().AddDate(0, 0, 6),
	})

	TestDummyNewsletters = append(TestDummyNewsletters, news.Newsletter{
		Id:                   bson.NewObjectId(),
		Beschreibung:         testBeschreibung + "2",
		BeschreibungEnglisch: testBeschreibungEnglisch + "2",
		Enddatum:             time.Now().AddDate(0, 0, 12),
		Person:               testPerson + "2",
		Startdatum:           time.Now(),
		Titel:                testTitel + "2",
		TitelEnglisch:        testTitelEnglisch + "2",
		Verdatum:             time.Now().AddDate(0, 0, 13),
	})

	TestDummyNewsletters = append(TestDummyNewsletters, news.Newsletter{
		Id:                   bson.NewObjectId(),
		Beschreibung:         testBeschreibung,
		BeschreibungEnglisch: testBeschreibungEnglisch,
		Enddatum:             time.Now().AddDate(0, 0, -5),
		Person:               testPerson,
		Startdatum:           time.Now(),
		Titel:                testTitel,
		TitelEnglisch:        testTitelEnglisch,
		Verdatum:             time.Now().AddDate(0, 0, 6),
	})
}

func TestNewsletterService(t *testing.T) {
	Setup()

	// setup mongodb session
	session, err := mongodb.NewSession(mongodbUrl)
	TestSession = *session
	if err != nil {
		t.Fatalf("Unable to connect to mongodb: %s", err)
	}
	// cleanup
	defer func() {
		// TestSession.DropDatabase(dbName)
		TestSession.Close()
	}()

	/* RUN-TEST: GetNewsletterTest
	 * INCLUDES TESTS:
	 * --- CreateNewsletters|GetNewsletterById : CreateNewsletterShouldInsertNewsletterIntoMongodb
	 * --- GetUpcomingNewsletters : GetUpcomingNewslettersShouldReturnOnlyNewslettersWithEndDateAfterNow
	 * --- DeleteNewsletterById : DeleteNewsletterShouldRemoveANewsletterFromMongodb
	 */
	t.Run("GetNewsletters", GetNewslettersShouldReturnMultipleNewslettersFromMongodb)
}

func GetNewslettersShouldReturnMultipleNewslettersFromMongodb(t *testing.T) {
	newsletterService := mongodb.NewNewsletterService(TestSession.Copy(), dbName, newsletterCollectionName)

	// add two entries to database
	t.Run("CreateNewsletters|GetNewsletterById", CreateNewsletterShouldInsertNewsletterIntoMongodb)
	err := newsletterService.CreateNewsletter(&TestDummyNewsletters[1])
	if err != nil {
		t.Errorf("Cannot create newsletter in database: %s", err)
	}

	// get newsletters
	newsletters, err := newsletterService.GetNewsletters()
	if err != nil {
		t.Fatal("Cannot get newsletters")
	}

	// check if atleast two entries
	count := len(newsletters)
	if count < 2 {
		t.Fatalf("Incorrect number of results. Expected: '2', got: '%d'", count)
	}
	t.Run("DeleteNewsletterById", DeleteNewsletterShouldRemoveANewsletterFromMongodb)
}

func CreateNewsletterShouldInsertNewsletterIntoMongodb(t *testing.T) {
	newsletterService := mongodb.NewNewsletterService(TestSession.Copy(), dbName, newsletterCollectionName)

	// store in database
	err := newsletterService.CreateNewsletter(&TestDummyNewsletters[0])
	if err != nil {
		t.Errorf("Cannot create newsletter in database: %s", err)
	}

	// test GetNewsletterById
	getNewsletter, err := newsletterService.GetNewsletterById(TestDummyNewsletters[0].Id.Hex())
	if err != nil {
		t.Fatalf("Cannot get newsletter with id: '%s", TestDummyNewsletters[0].Id)
	}

	if getNewsletter.Beschreibung != TestDummyNewsletters[0].Beschreibung &&
		getNewsletter.Verdatum != TestDummyNewsletters[0].Verdatum {
		t.Errorf("Result newsletter does not match created newsletter. "+
			"Expected: '%s' and '%s', "+
			"Got: '%s' and '%s'", TestDummyNewsletters[0].Beschreibung, TestDummyNewsletters[0].Verdatum.String(),
			getNewsletter.Beschreibung, getNewsletter.Verdatum.String())
	}
	log.Print(getNewsletter.Id)
}

func DeleteNewsletterShouldRemoveANewsletterFromMongodb(t *testing.T) {
	newsletterService := mongodb.NewNewsletterService(TestSession.Copy(), dbName, newsletterCollectionName)

	newsletters, err := newsletterService.GetNewsletters()
	if err != nil {
		t.Fatalf("Cannot create newsletter in database: %s", err)
	}

	count := len(newsletters)
	if count < 1 {
		t.Fatalf("No newsletters, test needs atleast one newsletter. Exected: 'atleast 1', Got: '%d'", count)
	}

	deleteNewsletter := newsletters[0]

	// delete newsletter
	err = newsletterService.DeleteNewsletterById(deleteNewsletter.Id.Hex())
	if err != nil {
		t.Fatalf("Cannot delete newsletter: '%s'", err)
	}

	// get new count after delete
	newCount, err := TestSession.GetCollection(dbName, newsletterCollectionName).Count()
	if err != nil {
		t.Errorf("Cannot query for count of newsletters in database: %s", err)
	}
	if newCount != count-1 {
		t.Errorf("Delete not sucessfull, to many entries remaining. Expected: '%d', Got: '%d'", count-1, newCount)
	}
}
