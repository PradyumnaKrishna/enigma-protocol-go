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
	User string `json:"user"`
}

type ConnectResponse struct {
	User      string `json:"user"`
	Publickey string `json:"publicKey"`
}

type TransmissionData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Payload string `json:"payload"`
}
