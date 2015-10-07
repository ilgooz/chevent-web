package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" schema:"name" validate:"required"`
	Date        time.Time     `json:"date" schema:"date" validate:"required"`
	Free        bool          `json:"free" schema:"free"`
	Image       string        `json:"image" schema:"image"`
	URL         string        `json:"url" schema:"url" validate:"required"`
	Description string        `json:"description" schema:"description"`
	Quota       int           `json:"quota" schema:"quota"`
	Speakers    []Speaker     `json:"speakers" schema:"speakers"`
}

type Speaker struct {
	Name    string `json:"name" schema:"name"`
	Subject string `json:"subject" schema:"subject"`
}
