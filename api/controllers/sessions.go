package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/reposandermets/tech-challenge-time/api/models"
)

type sessionController struct{}

type sessionControllerInterface interface {
	ListSession(w http.ResponseWriter, r *http.Request)
	PostSession(w http.ResponseWriter, r *http.Request)
	PutSessionByID(w http.ResponseWriter, r *http.Request)
	PutSessionBySessionID(w http.ResponseWriter, r *http.Request)
}

var (
	// SessionController ...
	SessionController sessionControllerInterface
)

func init() {
	SessionController = &sessionController{}
}

func (s *sessionController) ListSession(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "userID").(string)
	db := context.Get(r, "db").(*sql.DB)

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
	timeSessions, err = models.ReadOngoingSession(db, userID)
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
	userID := context.Get(r, "userID").(string)
	db := context.Get(r, "db").(*sql.DB)

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

	name := "My task started " + time.Now().Format(time.RFC850)
	if session.TimeSessionName.String != "" {
		name = session.TimeSessionName.String
	}
	sessionID, err := uuid.NewRandom()
	sessionIDString := sessionID.String()
	if session.TimeSessionID != "" {
		sessionIDString = session.TimeSessionID
		timeSessions, _ := models.ReadBySessionID(db, sessionIDString, userID)
		if len(timeSessions) < 1 {
			http.Error(w, http.StatusText(400), 400)
			log.Println(err.Error())
			return
		}
		name = timeSessions[0].TimeSessionName.String
	}

	timeSessions, err := models.WriteStartSession(db, sessionPartialID.String(), name, sessionIDString, userID)
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

func (s *sessionController) PutSessionByID(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "userID").(string)
	db := context.Get(r, "db").(*sql.DB)
	vars := mux.Vars(r)
	_, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println(err.Error())
		return
	}

	timeSessions, err := models.WriteStopSession(db, vars["id"], userID)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
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

func (s *sessionController) PutSessionBySessionID(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "userID").(string)
	db := context.Get(r, "db").(*sql.DB)
	vars := mux.Vars(r)
	_, err := uuid.Parse(vars["session_id"])
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		log.Println(err.Error())
		return
	}

	timeSessions, err := models.WriteEndSession(db, vars["session_id"], userID)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
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
