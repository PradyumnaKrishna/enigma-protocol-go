package db

import (
	"database/sql"
	"time"

	"enigma-protocol-go/pkg/models"
	"enigma-protocol-go/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	driver string
	uri    string
	conn   *sql.DB
}

type DatabaseOpts struct {
	Driver string
	Uri    string
}

func NewDatabase(dbopts DatabaseOpts) (*Database, error) {
	conn, err := sql.Open(dbopts.Driver, dbopts.Uri)
	if err != nil {
		return nil, err
	}

	db := &Database{
		driver: dbopts.Driver,
		uri:    dbopts.Uri,
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
		Driver: "sqlite3",
		Uri:    "users.db",
	})
}

func (d *Database) CreateTable() error {
	_, err := d.conn.Exec("CREATE TABLE IF NOT EXISTS Users (id TEXT PRIMARY KEY, publicKey TEXT, last_activity DATE)")
	if err != nil {
		return err
	}

	_, err = d.conn.Exec("CREATE TABLE IF NOT EXISTS PendingMessages (id INTEGER PRIMARY KEY AUTOINCREMENT, fromUser TEXT, toUser TEXT, payload TEXT)")
	return err
}

func (d *Database) GetPublicKey(id string) (string, error) {
	var key string
	err := d.conn.QueryRow("SELECT publicKey FROM Users WHERE id = ?", id).Scan(&key)
	return key, err
}

func (d *Database) SaveUser(publicKey string) (string, error) {
	id, err := utils.RandomHex(5)
	if err != nil {
		return "", err
	}

	stmt, err := d.conn.Prepare("INSERT INTO Users (id, publicKey, last_activity) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(id, publicKey, time.Now())
	return id, err
}

func (d *Database) IsUserExists(id string) bool {
	var count int
	d.conn.QueryRow("SELECT COUNT(*) FROM Users WHERE id = ?", id).Scan(&count)
	return count > 0
}

func (d *Database) UpdateActivity(id string) error {
	_, err := d.conn.Exec("UPDATE Users SET last_activity = ? WHERE id = ?", time.Now(), id)
	return err
}

func (d *Database) SavePendingMessage(message models.TransmissionData) error {
	stmt, err := d.conn.Prepare("INSERT INTO PendingMessages (fromUser, toUser, payload) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(message.From, message.To, message.Payload)
	return err
}

func (d *Database) GetPendingMessages(toUser string) ([]models.TransmissionData, error) {
	rows, err := d.conn.Query("SELECT fromUser, payload FROM PendingMessages WHERE toUser = ?", toUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.TransmissionData
	for rows.Next() {
		var fromUser, payload string
		err = rows.Scan(&fromUser, &payload)
		if err != nil {
			return nil, err
		}

		messages = append(messages, models.TransmissionData{
			From:    fromUser,
			To:      toUser,
			Payload: payload,
		})
	}

	return messages, nil
}

func (d *Database) DeletePendingMessages(toUser string) error {
	_, err := d.conn.Exec("DELETE FROM PendingMessages WHERE toUser = ?", toUser)
	return err
}
