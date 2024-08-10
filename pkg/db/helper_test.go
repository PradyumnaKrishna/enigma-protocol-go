package db

import (
	"enigma-protocol-go/pkg/models"
	"os"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		Driver: "sqlite3",
		Uri:    "test.db",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("Expected a Database instance, got nil")
	}
	defer db.conn.Close()
	defer os.Remove("test.db")
}

func TestNewDefaultDatabase(t *testing.T) {
	db, err := NewDefaultDatabase()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("Expected a Database instance, got nil")
	}
	defer db.conn.Close()
	defer os.Remove("users.db")
}

func TestCreateTable(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		Driver: "sqlite3",
		Uri:    "test.db",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.conn.Close()
	defer os.Remove("test.db")

	// Check if the table was created
	var tableName string
	err = db.conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", "Users").Scan(&tableName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSaveUser(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		Driver: "sqlite3",
		Uri:    "test.db",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.conn.Close()
	defer os.Remove("test.db")

	publicKey := "test-public-key"
	id, err := db.SaveUser(publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	key, err := db.GetPublicKey(id)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if key != publicKey {
		t.Fatalf("Expected public key %s, got %s", publicKey, key)
	}
}

func TestSavePendingMessage(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		Driver: "sqlite3",
		Uri:    "test.db",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.conn.Close()
	defer os.Remove("test.db")

	message := models.TransmissionData{
		From:    "test-from",
		To:      "test-to",
		Payload: "test-payload",
	}
	err = db.SavePendingMessage(message)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	messages, err := db.GetPendingMessages(message.To)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(messages) != 1 || messages[0] != message {
		t.Fatalf("Expected message %s, got %v", message, messages[0])
	}
}

func TestDeletePendingMessages(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		Driver: "sqlite3",
		Uri:    "test.db",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.conn.Close()
	defer os.Remove("test.db")

	message := models.TransmissionData{
		From:    "test-from",
		To:      "test-to",
		Payload: "test-payload",
	}
	err = db.SavePendingMessage(message)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = db.DeletePendingMessages(message.To)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	messages, err := db.GetPendingMessages(message.To)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(messages) != 0 {
		t.Fatalf("Expected no messages, got %v", messages)
	}
}
