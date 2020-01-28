package api

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/reposandermets/tech-challenge-time/api/config"
	"github.com/reposandermets/tech-challenge-time/api/controllers"
	"github.com/reposandermets/tech-challenge-time/api/models"
)

func loggingAndContext(db *sql.DB, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "db", db)
		log.Println(r.URL.Path)
		userID, err := uuid.Parse(r.Header.Get("x-user-uuid"))
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			log.Println("invalid header x-user-id")
			return
		}
		context.Set(r, "userID", userID.String())
		f(w, r)
	}
}

// Boot the service
func Boot() {
	config.SetEnvs()
	db, err := models.Connect()
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/v1/session-status", loggingAndContext(db, controllers.SessionController.ListSession)).Methods("GET")
	router.HandleFunc("/v1/session-start", loggingAndContext(db, controllers.SessionController.PostSession)).Methods("POST")
	router.HandleFunc("/v1/session-stop/{id}", loggingAndContext(db, controllers.SessionController.PutSessionByID)).Methods("PUT")
	router.HandleFunc("/v1/session-end/{session_id}", loggingAndContext(db, controllers.SessionController.PutSessionBySessionID)).Methods("PUT")

	server := &http.Server{
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "x-user-uuid"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router),
		Addr:         config.AppAddress,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
