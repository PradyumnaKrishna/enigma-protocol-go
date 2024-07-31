package api

import (
	"net/http"

	"enigma-protocol-go/pkg/db"
	"enigma-protocol-go/pkg/models"

	"github.com/julienschmidt/httprouter"
)

type ProtocolAPI struct {
	db *db.Database
}

func NewProtocolAPI(d db.Database) *ProtocolAPI {
	return &ProtocolAPI{db: &d}
}

func (p *ProtocolAPI) Register(r *httprouter.Router) {
	r.GET("/login/:publicKey", inJSON(p.login))
	r.GET("/connect/:id", inJSON(p.connect))
}

func (p *ProtocolAPI) login(_ *http.Request, ps httprouter.Params) (interface{}, *models.APIError) {
	publicKey := ps.ByName("publicKey")

	id, err := p.db.SaveUser(publicKey)
	if err != nil {
		return nil, &models.APIError{Code: http.StatusInternalServerError,
			Message: models.ErrorMessage{Error: "Internal Server Error", Detail: err.Error()},
		}
	}

	return &models.LoginResponse{ID: id}, nil
}

func (p *ProtocolAPI) connect(_ *http.Request, ps httprouter.Params) (interface{}, *models.APIError) {
	id := ps.ByName("id")

	publicKey, err := p.db.GetPublicKey(id)
	if err != nil {
		return nil, &models.APIError{Code: http.StatusNotFound,
			Message: models.ErrorMessage{Error: "Not Found"},
		}
	}

	return &models.ConnectResponse{User: id, Publickey: publicKey}, nil
}
