package commands

import (
	"os"
	"strings"
	"trash-rm/database"
	"trash-rm/parser"
	"trash-rm/utility"
)

func RestoreCommand(command parser.Command) error {
	// If only the Id was passed
	if command.Id > 0 {
		if err := restore(command); err != nil {
			return err
		}
	}
	return nil
}

// Restore by id
func restore(command parser.Command) error {
	// Get the informations from the database by id
	row, err := database.SelectById(command.Id)
	if err != nil {
		return err
	}

	// Check if the destination is a file
	isFile := strings.Contains(row.FromPath, ".")

	if isFile {
		// If destination is a file then uncompress the file
		if err := utility.UncompressFile(row.TrashPath, row.FromPath); err != nil {
			return err
		}
	} else {
		// Else uncompress the directory
		if err := utility.UncompressDir(row.TrashPath, row.FromPath); err != nil {
			return err
		}
	}
	
	// Delete unused entries in database
	if err := database.DeleteById(command.Id); err != nil {
		return err
	}

	// Delete compressed file
	if err := os.Remove(row.TrashPath); err != nil {
		return err
	}

	return nil
}