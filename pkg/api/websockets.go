package api

import (
	"context"
	"net/http"
	"sync"

	"enigma-protocol-go/pkg/db"
	"enigma-protocol-go/pkg/models"
	"nhooyr.io/websocket"
	"github.com/julienschmidt/httprouter"
)

type WebsocketAPI struct {
	db *db.Database
	connections map[string]*websocket.Conn
	mu sync.Mutex
}

func NewWebsocketAPI(d db.Database) *WebsocketAPI {
	return &WebsocketAPI{
		db: &d,
		connections: make(map[string]*websocket.Conn),
	}
}

func (w *WebsocketAPI) Register(r *httprouter.Router) {
	r.GET("/ws/:id", w.handleWebsocket)
}

func (w *WebsocketAPI) handleWebsocket(wr http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	conn, err := websocket.Accept(wr, r, nil)
	if err != nil {
		http.Error(wr, "Failed to establish websocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	w.mu.Lock()
	w.connections[id] = conn
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.connections, id)
		w.mu.Unlock()
	}()

	ctx := context.Background()

	pendingMessages, err := w.db.GetPendingMessages(id)
	if err != nil {
		http.Error(wr, "Failed to retrieve pending messages", http.StatusInternalServerError)
		return
	}

	for _, msg := range pendingMessages {
		err = conn.Write(ctx, websocket.MessageText, []byte(msg.Payload))
		if err != nil {
			http.Error(wr, "Failed to send pending message", http.StatusInternalServerError)
			return
		}
	}

	err = w.db.DeletePendingMessages(id)
	if err != nil {
		http.Error(wr, "Failed to delete pending messages", http.StatusInternalServerError)
		return
	}

	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			break
		}

		var message models.WebsocketMessage
		err = json.Unmarshal(msg, &message)
		if err != nil {
			http.Error(wr, "Invalid message format", http.StatusBadRequest)
			return
		}

		w.mu.Lock()
		receiverConn, connected := w.connections[message.To]
		w.mu.Unlock()

		if connected {
			err = receiverConn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				http.Error(wr, "Failed to send message", http.StatusInternalServerError)
				return
			}
		} else {
			err = w.db.SavePendingMessage(message.To, string(msg))
			if err != nil {
				http.Error(wr, "Failed to save pending message", http.StatusInternalServerError)
				return
			}
		}
	}
}
