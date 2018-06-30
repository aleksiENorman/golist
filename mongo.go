package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session = nil

func (e *entry) save() error {
	var err error
	var preExistingCount int
	if session == nil {
		session, err = mgo.Dial("localhost")
		if err != nil {
			return err
		}
	}

	c := session.DB("golist").C("posts")

	if e.ObjectId != "" {
		return c.UpdateId(e.ObjectId, e)
	}

	preExistingCount, err = c.Find(bson.M{"series": e.Series, "primary": true}).Count()
	if err != nil {
		return err
	}

	e.Primary = preExistingCount == 0
	e.ObjectId = bson.NewObjectId()
	return c.Insert(e)
}

func load(series string) ([]entry, error) {
	var err error
	var entries []entry

	if session == nil {
		session, err = mgo.Dial("localhost")
		if err != nil {
			return nil, err
		}
	}

	c := session.DB("golist").C("posts")
	if series != "index" {
		err = c.Find(bson.M{"series": series}).All(&entries)
	} else {
		err = c.Find(bson.M{"primary": true}).All(&entries)
	}

	return entries, err
}
