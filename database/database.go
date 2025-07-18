package database

import (
	"database/sql"
	"errors"
	"os"
	"path"

	_ "github.com/glebarez/go-sqlite"
)

// Database object
var Db *sql.DB

// Initialize database, create Db object and create tables if not exist
func InitDB() error {
	var err error
	// AppDataDir for linux
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("couldnt get home directory")
	}
	appDataDir := path.Join(homeDir, ".local/share/trash-rm")

	// Create or connect to database
	Db, err = sql.Open("sqlite", path.Join(appDataDir, "trm.db"))
	if err != nil {
		return errors.New("couldnt create or connect to database")
	}
	defer Db.Close()

	// Create tables if necessary
	if err := createTables(); err != nil {
		return err
	}
	
	return nil
}

// Create 2 tables
// Trash table were we store the actually data of the deleted files/folders
// And an extra tags table, because of database normalisation.
// One trash entry can have more then 1 tag. Basically a 1:n relation from trash to tags 
func createTables() error {
	// Trash table
	createTrashTable := "CREATE TABLE IF NOT EXISTS trash (" +
        "id 			INTEGER PRIMARY KEY," +
		"basename    	TEXT NOT NULL," +
        "from 			TEXT NOT NULL," +
		"to 			TEXT NOT NULL," +
		"createdAt   	DATETIME('now')"
	if _, err := Db.Exec(createTrashTable); err != nil {
		return errors.New("could not create trash table")
	}
	
	// Tags table
	createTagsTable := "CREATE TABLE IF NOT EXISTS tags (" +
		"id 			INTEGER PRIMARY KEY," +
		"tagname		TEXT NOT NULL," +
		"trash_id		INTEGER"
	if _, err := Db.Exec(createTagsTable); err != nil {
		return errors.New("could not create tags table")
	}
	
	return nil
}