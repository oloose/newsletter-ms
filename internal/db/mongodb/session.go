package mongodb

import (
	"errors"

	"github.com/globalsign/mgo"
)

// Defines a session type containing a session to mongodb.
type Session struct {
	session *mgo.Session
}

// Create a new session to mongodb (mHostPath is the path to mongodb).
func NewSession(mHostPath string) (*Session, error) {
	// try to get connection to mongodb and a session
	session, err := mgo.Dial(mHostPath)

	if err != nil {
		err = errors.New(err.Error() + "; Cannot connect to mongodb on " + mHostPath)
		return nil, err
	}

	return &Session{session}, err
}

// Returns a copy of the referenced session.
func (rSession *Session) Copy() *Session {
	return &Session{rSession.session.Copy()}
}

// Returns a collection from the referenced mongodb session.
func (rSession *Session) GetCollection(mDb string, mCollection string) *mgo.Collection {
	return rSession.session.DB(mDb).C(mCollection)
}

// Closes the mongodb session.
func (rSession *Session) Close() {
	if rSession.session != nil {
		rSession.session.Close()
	}
}

// Drops a given database (by its name mDb).
func (rSession *Session) DropDatabase(mDb string) error {
	if rSession.session != nil {
		return rSession.session.DB(mDb).DropDatabase()
	}

	return nil
}
