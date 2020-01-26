package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"

	"github.com/gorilla/mux"
	"github.com/reposandermets/tech-challenge-time/api/config"
	"github.com/reposandermets/tech-challenge-time/api/controllers"
	"github.com/reposandermets/tech-challenge-time/api/models"
)

func loggingAndContext(db *sql.DB, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "db", db)
		log.Println(r.URL.Path)
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

	router.HandleFunc("/v1/session", loggingAndContext(db, controllers.SessionController.ListSession)).Methods("GET")
	router.HandleFunc("/v1/session", loggingAndContext(db, controllers.SessionController.PostSession)).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         config.AppAddress,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
