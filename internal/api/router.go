package api

import (
	v1 "github.com/ejiro-edwin/todolist/internal/api/v1"
	"github.com/ejiro-edwin/todolist/internal/database"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//NewRouter provide a handler API service.
func NewRouter(db database.Database) (http.Handler, error) {
	router := mux.NewRouter()
	router.HandleFunc("/", v1.VersionHandler)
	router.HandleFunc("/version", v1.VersionHandler)

	v1.SetTodoAPI(db, router)

	router.Use(loggingMiddleware)

	return router, nil
}


func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		next.ServeHTTP(w, r)
		path := r.URL.Path
		end := time.Now().UTC()
		latency := end.Sub(start)
		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       path,
			"duration":   latency,
		}).Info()
	})
}
