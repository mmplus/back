package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
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
	a.Router.HandleFunc("/top", a.getTop).Methods("GET")
}

func (a *App) Run(addr string) {
	cors := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8888", handlers.CORS(cors)(a.Router)))
}

func (a *App) getRuns(w http.ResponseWriter, r *http.Request) {
	MongoRunCollection := a.MongoSession.DB("test").C("run")
	runs, err := getRuns(MongoRunCollection)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, runs)
}

func (a *App) getTop(w http.ResponseWriter, r *http.Request) {
	MongoRunCollection := a.MongoSession.DB("test").C("run")
	runs, err := getTop(MongoRunCollection)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, runs)
}
