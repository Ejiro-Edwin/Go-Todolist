package main

import (
	"flag"
	"github.com/ejiro-edwin/todolist/internal/database"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ejiro-edwin/todolist/internal/api"
	"github.com/ejiro-edwin/todolist/internal/config"
)

func main(){
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithField("version", config.Version).Debug("Starting server.")

	//Creating new database
	db, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Error verifying database.")
	}

	logrus.Debug("Database is ready to use.")

	//Creating new router
	router, err := api.NewRouter(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8088"
	}

	var addr = "0.0.0.0:" + port
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	//Starting server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed.")
	}
}

