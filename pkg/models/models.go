package models

type APIError struct {
	Code    int          `json:"code"`
	Message ErrorMessage `json:"message"`
}

type ErrorMessage struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

type LoginResponse struct {
	ID string `json:"id"`
}

type ConnectResponse struct {
	User      string `json:"user"`
	Publickey string `json:"publickey"`
}

type TransmissionData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Payload string `json:"payload"`
}

type WebsocketMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Payload string `json:"payload"`
}

type PendingMessage struct {
	To      string `json:"to"`
	Payload string `json:"payload"`
}
