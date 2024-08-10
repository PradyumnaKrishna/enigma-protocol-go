package main

import (
	"enigma-protocol-go/pkg/api"
	"enigma-protocol-go/pkg/db"
	"log"
	"net/http"
	"os"
	"strings"
)

type opts struct {
	port           string
	databasePath   string
	allowedOrigins []string
}

func main() {
	env := getEnv()
	apiOpts, err := api.NewAPIOpts(
		&db.DatabaseOpts{
			Driver: "sqlite3",
			Uri:    env.databasePath,
		},
		env.allowedOrigins,
	)
	if err != nil {
		panic(err)
	}

	router := apiOpts.NewRouter()

	log.Println("Server Configuration")
	log.Printf("Port: %s\n", env.port)
	log.Printf("Database Path: %s\n", env.databasePath)
	log.Printf("Allowed Origins: %s\n", env.allowedOrigins)

	log.Println("Starting server on :" + env.port)
	http.ListenAndServe(":"+env.port, router)
}

func getEnv() opts {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "sqlite3.db"
	}

	return opts{
		port:           port,
		databasePath:   databasePath,
		allowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
	}
}
