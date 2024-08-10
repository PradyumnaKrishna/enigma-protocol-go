package api

import (
	"encoding/json"
	"net/http"

	"enigma-protocol-go/pkg/db"
	"enigma-protocol-go/pkg/models"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type APIFunc func(r *http.Request, ps httprouter.Params) (interface{}, *models.APIError)

type APIOpts struct {
	Database       *db.Database
	AllowedOrigins []string
}

func NewAPIOpts(
	dbopts *db.DatabaseOpts,
	allowedOrigins []string,
) (*APIOpts, error) {
	var database *db.Database
	var err error

	if dbopts == nil {
		database, err = db.NewDefaultDatabase()
	} else {
		database, err = db.NewDatabase(*dbopts)
	}

	if err != nil {
		return nil, err
	}

	return &APIOpts{
		Database:       database,
		AllowedOrigins: allowedOrigins,
	}, nil
}

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

func (opts APIOpts) NewRouter() http.Handler {
	router := httprouter.New()

	protocolAPI := NewProtocolAPI(opts)
	protocolAPI.Register(router)

	websocketAPI := NewWebsocketAPI(opts)
	websocketAPI.Register(router)

	router.GET("/", inJSON(index))
	router.GET("/version", inJSON(version))

	_cors := cors.Options{
		AllowedOrigins: opts.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST"},
	}

	handler := cors.New(_cors).Handler(router)
	return handler
}

func index(r *http.Request, ps httprouter.Params) (interface{}, *models.APIError) {
	return map[string]string{"status": "ok"}, nil
}

func version(r *http.Request, ps httprouter.Params) (interface{}, *models.APIError) {
	return map[string]string{"version": "0.3.0"}, nil
}
