package db

import (
	"database/sql"
	"fmt"
	"time"

	"enigma-protocol-go/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	uri    string
	table  string
	driver string
	conn   *sql.DB
}

type DatabaseOpts struct {
	uri    string
	table  string
	driver string
}

func NewDatabase(dbopts DatabaseOpts) (*Database, error) {
	conn, err := sql.Open(dbopts.driver, dbopts.uri)
	if err != nil {
		return nil, err
	}

	db := &Database{
		uri:    dbopts.uri,
		table:  dbopts.table,
		driver: dbopts.driver,
		conn:   conn,
	}

	err = db.CreateTable()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewDefaultDatabase() (*Database, error) {
	return NewDatabase(DatabaseOpts{
		uri:    "users.db",
		table:  "Users",
		driver: "sqlite3",
	})
}

func (d *Database) CreateTable() error {
	_, err := d.conn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id TEXT PRIMARY KEY, publicKey TEXT, last_activity DATE)", d.table))
	return err
}

func (d *Database) GetPublicKey(id string) (string, error) {
	var key string
	err := d.conn.QueryRow(fmt.Sprintf("SELECT publicKey FROM %s WHERE id = ?", d.table), id).Scan(&key)
	return key, err
}

func (d *Database) SaveUser(publicKey string) (string, error) {
	id, err := utils.RandomHex(5)
	if err != nil {
		return "", err
	}

	stmt, err := d.conn.Prepare(fmt.Sprintf("INSERT INTO %s (id, publicKey, last_activity) VALUES (?, ?, ?)", d.table))
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(id, publicKey, time.Now())
	return id, err
}

func (d *Database) IsUserExists(id string) bool {
	var count int
	d.conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", d.table), id).Scan(&count)
	return count > 0
}

func (d *Database) UpdateActivity(id string) error {
	_, err := d.conn.Exec(fmt.Sprintf("UPDATE %s SET last_activity = ? WHERE id = ?", d.table), time.Now(), id)
	return err
}
