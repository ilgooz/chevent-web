package conf

import (
	"log"

	"gopkg.in/mgo.v2"
)

func dialMongo() *mgo.Session {
	session, err := mgo.Dial(*mongoAddr)

	if err != nil {
		log.Fatalln(err)
	}

	session.SetSafe(&mgo.Safe{})

	return session
}
