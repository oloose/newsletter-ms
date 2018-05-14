package mongodb_test

import (
	"testing"
	"newsletter-service/internal/news"
	"time"
	"newsletter-service/internal/db/mongodb"
	"log"
)

var TestSession mongodb.Session
var TestDummyNewsletters []news.Newsletter

const(
	mongodbUrl = "localhost:27017"
	dbName = "test_newsletterdb"
	newsletterCollectionName = "test_NewsletterCollection"

	testBeschreibung = "testBeschreibung"
	testBeschreibungEnglisch = "testBeschreibungEnglisch"
	testPerson = "testPerson"
	testTitel = "testTitel"
	testTitelEnglisch = "testTitelEnglisch"
)

func Setup(){

	// setup dummy newsletters
	TestDummyNewsletters = append(TestDummyNewsletters, news.Newsletter{
		Beschreibung: testBeschreibung,
		BeschreibungEnglisch: testBeschreibungEnglisch,
		Enddatum: time.Now().AddDate(0, 0, 5),
		Person: testPerson,
		Startdatum: time.Now(),
		Titel: testTitel,
		TitelEnglisch: testTitelEnglisch,
		Verdatum: time.Now().AddDate(0,0,6),
	})

	TestDummyNewsletters = append(TestDummyNewsletters, news.Newsletter{
		Beschreibung: testBeschreibung + "2",
		BeschreibungEnglisch: testBeschreibungEnglisch + "2",
		Enddatum: time.Now().AddDate(0, 0, 12),
		Person: testPerson + "2",
		Startdatum: time.Now(),
		Titel: testTitel + "2",
		TitelEnglisch: testTitelEnglisch + "2",
		Verdatum: time.Now().AddDate(0,0,13),
	})

	TestDummyNewsletters = append(TestDummyNewsletters, news.Newsletter{
		Beschreibung: testBeschreibung,
		BeschreibungEnglisch: testBeschreibungEnglisch,
		Enddatum: time.Now().AddDate(0, 0, -5),
		Person: testPerson,
		Startdatum: time.Now(),
		Titel: testTitel,
		TitelEnglisch: testTitelEnglisch,
		Verdatum: time.Now().AddDate(0,0,6),
	})
}

func TestNewsletterService(t *testing.T){
	Setup()

	// setup mongodb session
	session, err := mongodb.NewSession(mongodbUrl)
	TestSession = *session
	if err != nil{
		t.Fatalf("Unable to connect to mongodb: %s", err)
	}
	// cleanup
	defer func() {
		TestSession.DropDatabase(dbName)
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

func GetNewslettersShouldReturnMultipleNewslettersFromMongodb(t *testing.T){
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
	if count < 2{
		t.Fatalf("Incorrect number of results. Expected: '2', got: '%d'", count)
	}
	t.Run("GetUpcomingNewsletters", GetUpcomingNewslettersShouldReturnOnlyNewslettersWithEndDateAfterNow)
	t.Run("DeleteNewsletterById", DeleteNewsletterShouldRemoveANewsletterFromMongodb)
}

func CreateNewsletterShouldInsertNewsletterIntoMongodb(t *testing.T){
	newsletterService := mongodb.NewNewsletterService(TestSession.Copy(), dbName, newsletterCollectionName)

	// store in database
	err := newsletterService.CreateNewsletter(&TestDummyNewsletters[0])
	if err != nil {
		t.Errorf("Cannot create newsletter in database: %s", err)
	}

	// Check if newsletter was created
	var results []news.Newsletter
	err = TestSession.GetCollection(dbName, newsletterCollectionName).Find(nil).All(&results)
	if err != nil  {
		t.Fatalf("Unable to get newsletters from database: '%s'", err)
	}
	if len(results) < 1{
		t.Fatalf("No newsletter entires found in database. Exected: 'atleast 1', Got: '%d'", len(results))
	}

	// test GetNewsletterById
	getNewsletter, err := newsletterService.GetNewsletterById(results[0].Id.Hex())
	if err != nil {
		t.Fatalf("Cannot get newsletter with id: '%s", results[0].Id)
	}
	count := len(results)
	// only one entry should be found
	if count != 1{
		t.Errorf("Incorrect number of results. Expected: '1', got: %d", count)
	}
	if getNewsletter.Beschreibung != TestDummyNewsletters[0].Beschreibung &&
		getNewsletter.Verdatum != TestDummyNewsletters[0].Verdatum{
			t.Errorf("Result newsletter does not match created newsletter. " +
				"Expected: '%s' and '%s', " +
				"Got: '%s' and '%s'", TestDummyNewsletters[0].Beschreibung, TestDummyNewsletters[0].Verdatum.String(),
				getNewsletter.Beschreibung, getNewsletter.Verdatum.String())
	}
	log.Print(getNewsletter.Id)
}

func DeleteNewsletterShouldRemoveANewsletterFromMongodb(t *testing.T){
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
	if newCount != count - 1{
		t.Errorf("Delete not sucessfull, to many entries remaining. Expected: '%d', Got: '%d'", count - 1, newCount)
	}
}

func GetUpcomingNewslettersShouldReturnOnlyNewslettersWithEndDateAfterNow(t *testing.T){

}

