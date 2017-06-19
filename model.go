package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Player struct {
	Name   string
	Armory string
	Role   string
}

type Metadata struct {
	Patch string
	Affix []string
}

type RunMetadata struct {
	Realm   string
	Region  string
	Dungeon string
}

type Run struct {
	Id          string `json:"id" bson:"_id,omitempty"`
	Level       int64
	Time        string
	Party       []Player
	Completed   time.Time
	RunMetadata RunMetadata
	Metadata    Metadata
}

func (r *Run) getRun(collection *mgo.Collection) error {
	return errors.New("Not implemented")
}

func (r *Run) createRun(collection *mgo.Collection) error {
	return errors.New("Not implemented")
}

func getRuns(collection *mgo.Collection) ([]Run, error) {
	runs := []Run{}
	err := collection.Find(bson.M{}).All(&runs)
	if err != nil {
		return runs, errors.New("Not implemented")
	}
	return runs, nil
}
