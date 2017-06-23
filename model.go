package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Player struct {
	Name   string `json:"name"`
	Armory string `json:"armory"`
	Role   string `json:"role"`
}

type Metadata struct {
	Patch string   `json:"patch"`
	Affix []string `json:"affix"`
}

type RunMetadata struct {
	Realm   string `json:"realm"`
	Region  string `json:"region"`
	Dungeon string `json:"dungeon"`
}

type Run struct {
	Id          string      `json:"id" bson:"_id,omitempty"`
	Level       int64       `json:"level"`
	Time        string      `json:"time"`
	Party       []Player    `json:"party"`
	Completed   time.Time   `json:"completed"`
	RunMetadata RunMetadata `json:"run_metadata"`
	Metadata    Metadata    `json:"metadata"`
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

func getTop(collection *mgo.Collection) ([]Run, error) {
	runs := []Run{}
	err := collection.Find(bson.M{}).Limit(100).Sort("-level").All(&runs)
	if err != nil {
		return runs, errors.New("Not implemented")
	}
	return runs, nil
}
