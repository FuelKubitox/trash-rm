package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"trash-rm/database"
	"trash-rm/parser"
)

// Empty the trash
func EmptyCommand(command parser.Command) error {
	var trashList []database.TrashRow
	var err error
	// Get all entries from the databse
	result, err := database.SelectAll()
	if err != nil {
		return err
	}
	// Parse the data
	trashList, err = database.ParseDbData(result)
	if err != nil {
		return err
	}
	// Ask the user if he wants to delete everything
	r := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "Trash has " + strconv.Itoa(len(trashList)) + " entries. Are you sure you want to delete them all? [y/n] ")
	str, _ := r.ReadString('\n')
    action := strings.TrimSpace(str)
	if action == "y" {
		if err := emptyTrash(trashList); err != nil {
			return err
		}
	}

	return nil
}

// Loop through all entries and delete every single file
func emptyTrash(trashList []database.TrashRow) error {
	// First delete the files
	for _, row := range trashList {
		fmt.Println("Deleting " + row.TrashPath)
		if err := os.Remove(row.TrashPath); err != nil {
			fmt.Println("Could not delete " + row.TrashPath)
		}
	}
	// Then drop the database tables
	if err := database.DropTables(); err != nil {
		return err
	}
	fmt.Println("Trash deleted!")

	return nil
}