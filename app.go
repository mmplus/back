package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type App struct {
	Router       *mux.Router
	MongoSession *mgo.Session
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) Initialize(username, password, dburi string) {
	var err error

	a.MongoSession, err = mgo.Dial("mongodb://localhost")

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/runs", a.getRuns).Methods("GET")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8888", a.Router))
}

func (a *App) getRuns(w http.ResponseWriter, r *http.Request) {
	MongoRunCollection := a.MongoSession.DB("mmplus").C("run")
	runs, err := getRuns(MongoRunCollection)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, runs)
}
