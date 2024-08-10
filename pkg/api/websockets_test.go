package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"enigma-protocol-go/pkg/models"

	"nhooyr.io/websocket"
)

func createUser(t *testing.T, router http.Handler, publicKey string) string {
	req, _ := http.NewRequest("GET", "/login/"+publicKey, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
	}

	var res models.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	return res.User
}

func TestConnectInvalidUser(t *testing.T) {
	router := setup()
	defer cleanup()

	s := httptest.NewServer(router)
	wsEndpoint := "ws" + s.URL[4:] + "/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, wsEndpoint+"random-user", nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	_, msg, err := c.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var error models.ErrorMessage
	err = json.Unmarshal(msg, &error)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if error.Error != "User not found" {
		t.Errorf("Expected User not found, got %v", error.Error)
	}
}

func TestSendToInvalidUser(t *testing.T) {
	router := setup()
	defer cleanup()

	user1 := createUser(t, router, "key1")

	s := httptest.NewServer(router)
	wsEndpoint := "ws" + s.URL[4:] + "/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, wsEndpoint+user1, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	data, _ := json.Marshal(models.TransmissionData{
		From:    user1,
		To:      "random-user-new",
		Payload: "Hello User",
	})
	err = c.Write(ctx, websocket.MessageText, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, msg, err := c.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var error models.ErrorMessage
	err = json.Unmarshal(msg, &error)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if error.Error != "User not found" {
		t.Errorf("Expected User not found, got %v", error.Error)
	}
}

func TestSyncCommunication(t *testing.T) {
	router := setup()
	defer cleanup()

	user1 := createUser(t, router, "key1")
	user2 := createUser(t, router, "key2")

	s := httptest.NewServer(router)
	wsEndpoint := "ws" + s.URL[4:] + "/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Connect user1
	c1, _, err := websocket.Dial(ctx, wsEndpoint+user1, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	// Connect user2
	c2, _, err := websocket.Dial(ctx, wsEndpoint+user2, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	// Send message from user1 to user2
	data := models.TransmissionData{
		From:    user1,
		To:      user2,
		Payload: "Hello User",
	}
	jsonData, _ := json.Marshal(data)
	err = c1.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read message from user1 using user2 connection
	_, msg, err := c2.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var response1 models.TransmissionData
	err = json.Unmarshal(msg, &response1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the message is received correctly
	if !reflect.DeepEqual(response1, data) {
		t.Errorf("Expected %v, got %v", data, response1)
	}

	// Send message from user2 to user1
	data = models.TransmissionData{
		From:    user2,
		To:      user1,
		Payload: "Hello Another User",
	}
	jsonData, _ = json.Marshal(data)
	err = c2.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read message from user2 using user1 connection
	_, msg, err = c1.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var response2 models.TransmissionData
	err = json.Unmarshal(msg, &response2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the message is received correctly
	if !reflect.DeepEqual(response2, data) {
		t.Errorf("Expected %v, got %v", data, response2)
	}
}

func TestAsyncCommunication(t *testing.T) {
	router := setup()
	defer cleanup()

	user1 := createUser(t, router, "key1")
	user2 := createUser(t, router, "key2")

	s := httptest.NewServer(router)
	wsEndpoint := "ws" + s.URL[4:] + "/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c1, _, err := websocket.Dial(ctx, wsEndpoint+user1, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	c2, _, err := websocket.Dial(ctx, wsEndpoint+user2, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	data := models.TransmissionData{
		From:    user1,
		To:      user2,
		Payload: "Hello User",
	}
	jsonData, _ := json.Marshal(data)
	err = c1.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, msg, err := c2.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var response1 models.TransmissionData
	err = json.Unmarshal(msg, &response1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(response1, data) {
		t.Errorf("Expected %v, got %v", data, response1)
	}

	data = models.TransmissionData{
		From:    user2,
		To:      user1,
		Payload: "Hello Another User",
	}
	jsonData, _ = json.Marshal(data)
	err = c2.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, msg, err = c1.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var response2 models.TransmissionData
	err = json.Unmarshal(msg, &response2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(response2, data) {
		t.Errorf("Expected %v, got %v", data, response2)
	}
}

func TestPendingMessages(t *testing.T) {
	router := setup()
	defer cleanup()

	user1 := createUser(t, router, "key1")
	user2 := createUser(t, router, "key2")

	s := httptest.NewServer(router)
	wsEndpoint := "ws" + s.URL[4:] + "/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c1, _, err := websocket.Dial(ctx, wsEndpoint+user1, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	data := models.TransmissionData{
		From:    user1,
		To:      user2,
		Payload: "Hello User",
	}
	jsonData, _ := json.Marshal(data)
	err = c1.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// time.Sleep(100 * time.Millisecond)
	c2, _, err := websocket.Dial(ctx, wsEndpoint+user2, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	_, msg, err := c2.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var response1 models.TransmissionData
	err = json.Unmarshal(msg, &response1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(response1, data) {
		t.Errorf("Expected %v, got %v", data, response1)
	}

	data = models.TransmissionData{
		From:    user2,
		To:      user1,
		Payload: "Hello Another User",
	}
	jsonData, _ = json.Marshal(data)
	err = c2.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	c2.Close(websocket.StatusNormalClosure, "")

	_, msg, err = c1.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var response2 models.TransmissionData
	err = json.Unmarshal(msg, &response2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(response2, data) {
		t.Errorf("Expected %v, got %v", data, response2)
	}

	data = models.TransmissionData{
		From:    user1,
		To:      user2,
		Payload: fmt.Sprintf("Hello User %s", user2),
	}
	jsonData, _ = json.Marshal(data)
	err = c1.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	c1.Close(websocket.StatusNormalClosure, "")

	c2, _, err = websocket.Dial(ctx, wsEndpoint+user2, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	_, msg, err = c2.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = json.Unmarshal(msg, &response1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(response1, data) {
		t.Errorf("Expected %v, got %v", data, response1)
	}
}

func TestSendInvalidData(t *testing.T) {
	router := setup()
	defer cleanup()

	user1 := createUser(t, router, "key1")

	s := httptest.NewServer(router)
	wsEndpoint := "ws" + s.URL[4:] + "/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c1, _, err := websocket.Dial(ctx, wsEndpoint+user1, nil)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	err = c1.Write(ctx, websocket.MessageText, []byte("invalid"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, msg, err := c1.Read(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var error models.ErrorMessage
	err = json.Unmarshal(msg, &error)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if error.Error != "Invalid message format" {
		t.Errorf("Expected Invalid data, got %v", error.Error)
	}
}
