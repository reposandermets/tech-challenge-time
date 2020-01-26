package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/reposandermets/tech-challenge-time/api/models"
)

type sessionController struct{}

type sessionControllerInterface interface {
	ListSession(w http.ResponseWriter, r *http.Request)
	PostSession(w http.ResponseWriter, r *http.Request)
}

var (
	// SessionController ...
	SessionController sessionControllerInterface
)

func init() {
	SessionController = &sessionController{}
}

func (s *sessionController) ListSession(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.Header.Get("x-user-uuid"))
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println("invalid header x-user-id")
		return
	}

	dbRaw, ok := context.GetOk(r, "db")
	db := dbRaw.(*sql.DB)
	if !ok {
		log.Println("could not get database connection pool from context")
		http.Error(w, "please retry", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// active tracking
	timeSessions, err := models.ReadOngoingTracking(db, userID.String())
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println("db read error", err.Error())
		return
	}

	if len(timeSessions) > 0 {
		b, err := json.Marshal(timeSessions)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			log.Println("marshal error", err.Error())
			return
		}
		w.WriteHeader(200)
		w.Write(b)
		return
	}

	// ongoing session
	timeSessions, err = models.ReadOngoingSession(db, userID.String())
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println("db read error", err.Error())
		return
	}

	if len(timeSessions) > 0 {
		b, err := json.Marshal(timeSessions)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			log.Println("marshal error", err.Error())
			return
		}
		w.WriteHeader(200)
		w.Write(b)
		return
	}

	// no active sessions nor ongoing tracking
	w.WriteHeader(200)
	w.Write([]byte(`[]`))
	return
}

func (s *sessionController) PostSession(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.Header.Get("x-user-uuid"))
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println("invalid header x-user-id")
		return
	}

	dbRaw, ok := context.GetOk(r, "db")
	db := dbRaw.(*sql.DB)
	if !ok {
		log.Println("could not get database connection pool from context")
		http.Error(w, "please retry", 500)
		return
	}

	var session models.TimeSession
	sessionBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println(err.Error())
		return
	}

	if err := json.Unmarshal(sessionBytes, &session); err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println(err.Error())
		return
	}
	sessionPartialID, err := uuid.NewRandom()

	name := "My task"
	if session.TimeSessionName.String != "" {
		name = session.TimeSessionName.String
	}
	sessionID, err := uuid.NewRandom()
	sessionIDString := sessionID.String()
	if session.TimeSessionID != "" {
		sessionIDString = session.TimeSessionID
		timeSessions, _ := models.ReadBySessionID(db, sessionIDString, userID.String())
		if len(timeSessions) < 1 {
			http.Error(w, http.StatusText(400), 400)
			log.Println(err.Error())
			return
		}
	}

	// no active sessions nor ongoing tracking
	timeSessions, err := models.WriteStartSession(db, sessionPartialID.String(), name, sessionIDString, userID.String())
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Println(err.Error())
		return
	}
	b, err := json.Marshal(timeSessions)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("marshal error", err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(b)
	return
}
