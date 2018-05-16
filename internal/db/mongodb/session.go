package mongodb

import (
	"github.com/globalsign/mgo"
)

type Session struct{
	session *mgo.Session
}

func NewSession(mUrl string) (*Session, error){
	// "localhost:27017 is base path to mongodb
	session, err := mgo.Dial(mUrl)

	if err != nil{
		return nil, err
	}

	return &Session{session}, err
}

func(rSession *Session) Copy() *Session{
	return &Session{rSession.session.Copy()}
}

func(rSession *Session) GetCollection(mDb string, mCollection string) *mgo.Collection{
	return rSession.session.DB(mDb).C(mCollection)
}

func(rSession *Session) Close(){
	if rSession.session != nil {
		rSession.session.Close()
	}
}

func(rSession *Session) DropDatabase(mDb string) error{
	if rSession.session != nil{
		return rSession.session.DB(mDb).DropDatabase()
	}

	return nil
}
