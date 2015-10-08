package database

import "gopkg.in/mgo.v2"

type Mongo struct {
	session *mgo.Session
	mongo   *mgo.Database
}

func New() *Mongo {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return &Mongo{session, session.DB("db")}
}

func (db *Mongo) Save(collectionName string, data interface{}) error {
	return db.mongo.C(collectionName).Insert(data)
}

func (db *Mongo) Close() {
	db.session.Close()
}
