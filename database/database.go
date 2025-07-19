package database

import (
	"database/sql"
	"errors"
	"fmt"
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
		fmt.Println(err)
		return errors.New("couldnt get home directory")
	}
	appDataDir := path.Join(homeDir, ".local/share/trash-rm")
	if err = os.MkdirAll(appDataDir, os.ModePerm); err != nil {
		fmt.Println(err)
		return errors.New("couldnt create app data directory")
	}

	// Create or connect to database
	Db, err = sql.Open("sqlite", path.Join(appDataDir, "trm.db"))
	if err != nil {
		fmt.Println(err)
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
	createTrashTable := "CREATE TABLE IF NOT EXISTS trash_table (" +
        "trash_id INTEGER PRIMARY KEY," +
		"basename TEXT NOT NULL," +
        "from_path TEXT NOT NULL," +
		"trash_path TEXT NOT NULL," +
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP);"
	if _, err := Db.Exec(createTrashTable); err != nil {
		fmt.Println(err)
		return errors.New("could not create trash table in database")
	}
	
	// Tags table
	createTagsTable := "CREATE TABLE IF NOT EXISTS tags_table (" +
		"tags_id INTEGER PRIMARY KEY," +
		"tagname TEXT NOT NULL," +
		"trash_tag_id INTEGER NOT NULL," +
		"FOREIGN KEY(trash_tag_id) REFERENCES trash_table(trash_id));"
	if _, err := Db.Exec(createTagsTable); err != nil {
		fmt.Println(err)
		return errors.New("could not create tags table in database")
	}
	
	return nil
}

// Insert a new trash entry without tags
func Insert(basename string, from string, to string) error {
	insert := "INSERT INTO trash_table " +
		"(basename, from_path, trash_path) VALUES (" +
		basename + "," +
		from + "," +
		to + ");"
	if _, err := Db.Exec(insert); err != nil {
		fmt.Println(err)
		return errors.New("could not create trash entry in database")
	}

	return nil
}