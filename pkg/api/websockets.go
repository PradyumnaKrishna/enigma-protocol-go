package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"enigma-protocol-go/pkg/db"
	"enigma-protocol-go/pkg/models"

	"github.com/julienschmidt/httprouter"
	"nhooyr.io/websocket"
)

type WebsocketAPI struct {
	db    *db.Database
	chats map[string]Chat
	mu    sync.Mutex
}

func NewWebsocketAPI(opts APIOpts) *WebsocketAPI {
	return &WebsocketAPI{
		db:    opts.Database,
		chats: make(map[string]Chat),
	}
}

type Chat struct {
	connection *websocket.Conn
}

func (chat *Chat) sendJSON(ctx context.Context, message interface{}) error {
	if chat.connection == nil {
		return nil
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return chat.connection.Write(ctx, websocket.MessageText, []byte(data))
}

func (chat *Chat) SendMessage(ctx context.Context, message models.TransmissionData) error {
	return chat.sendJSON(ctx, message)
}

func (chat *Chat) sendPendingMessages(messages []models.TransmissionData) error {
	for _, message := range messages {
		err := chat.SendMessage(context.Background(), message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WebsocketAPI) Register(r *httprouter.Router) {
	r.GET("/ws/:id", w.handleWebsocket)
}

func (w *WebsocketAPI) handleWebsocket(wr http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	conn, err := websocket.Accept(wr, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})

	if err != nil {
		http.Error(wr, "Failed to establish websocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close(websocket.StatusInvalidFramePayloadData, "Internal Error")

	ctx := context.Background()
	chat := Chat{connection: conn}
	if !w.db.IsUserExists(id) {
		chat.sendJSON(ctx, models.ErrorMessage{
			Error: "User not found",
		})
		return
	}

	// if user already connected, close the connection
	w.mu.Lock()
	if _, ok := w.chats[id]; ok {
		chat.sendJSON(ctx, models.ErrorMessage{
			Error: "User connected from another location",
		})
		return
	} else {
		w.chats[id] = chat
	}
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.chats, id)
		w.mu.Unlock()
	}()

	pendingMessages, _ := w.db.GetPendingMessages(id)
	err = chat.sendPendingMessages(pendingMessages)

	if err == nil {
		w.db.DeletePendingMessages(id)
	}

	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			break
		}

		var message models.TransmissionData
		err = json.Unmarshal(msg, &message)
		if err != nil {
			chat.sendJSON(ctx, models.ErrorMessage{
				Error: "Invalid message format",
			})
			continue
		}

		w.mu.Lock()
		receiverConn, connected := w.chats[message.To]
		w.mu.Unlock()

		if connected {
			receiverConn.SendMessage(ctx, message)
		} else {
			if w.db.IsUserExists(message.To) {
				w.db.SavePendingMessage(message)
			} else {
				chat.sendJSON(ctx, models.ErrorMessage{
					Error: "User not found",
				})
			}
		}
	}
}
