package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"enigma-protocol-go/pkg/db"
	"enigma-protocol-go/pkg/models"

	"github.com/julienschmidt/httprouter"
)

type APIFunc func(r *http.Request, ps httprouter.Params) (interface{}, *models.APIError)

func inJSON(api APIFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		res, err := api(r, ps)
		if err != nil {
			w.WriteHeader(err.Code)
			json.NewEncoder(w).Encode(err.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func StartServer() {
	router := httprouter.New()

	database, err := db.NewDefaultDatabase()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	protocolAPI := NewProtocolAPI(*database)
	protocolAPI.Register(router)

	router.GET("/", inJSON(index))
	router.GET("/version", inJSON(version))

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}

func index(r *http.Request, ps httprouter.Params) (interface{}, *models.APIError) {
	return map[string]string{"status": "ok"}, nil
}

func version(r *http.Request, ps httprouter.Params) (interface{}, *models.APIError) {
	return map[string]string{"version": "0.3.0"}, nil
}
