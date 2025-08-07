package commands

import (
	"fmt"
	"os"

	"trash-rm/database"
	"trash-rm/parser"
)

// Synchronize the data in the database with the files on the filesystem
// If there is a database entry with no existing file in the filesystem,
// that database entry gets deleted
func SyncCommand(command parser.Command) error {
	// Get trash database data
	result, err := database.SelectAll()
	if err != nil {
		return err
	}

	// Parse trash database data
	trashList, err := database.ParseDbData(result)
	if err != nil {
		return err
	}

	// Start syncing database with files on the filesystem
	if err := sync(trashList); err != nil {
		return err
	}

	return nil
}

// Start synchronize
func sync(trashList []database.TrashRow) error {
	fmt.Println("Start deleting database entries pointing to files that doesnt exist...")
	// Loop through all database entries and check if the file exists
	for _, row := range trashList {
		_, err := os.Stat(row.TrashPath)
		// If the file does not exist delete the database entry
		if os.IsNotExist(err) {
			err := database.DeleteById(row.Id)
			if err != nil {
				fmt.Println("Couldnt delete database entry during sync")
				return err
			}
			fmt.Println("Deleted database entry with filepath " + row.TrashPath)
		}
	}
	fmt.Println("Finished syncing!")

	return nil
}