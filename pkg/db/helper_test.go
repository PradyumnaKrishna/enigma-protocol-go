package db

import (
	"os"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		uri:    "test.db",
		table:  "TestTable",
		driver: "sqlite3",
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
		uri:    "test.db",
		table:  "TestTable",
		driver: "sqlite3",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.conn.Close()
	defer os.Remove("test.db")

	// Check if the table was created
	var tableName string
	err = db.conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", db.table).Scan(&tableName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if tableName != db.table {
		t.Fatalf("Expected table name %s, got %s", db.table, tableName)
	}
}

func TestSaveUser(t *testing.T) {
	db, err := NewDatabase(DatabaseOpts{
		uri:    "test.db",
		table:  "TestTable",
		driver: "sqlite3",
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
