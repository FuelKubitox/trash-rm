package commands

import (
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

	// Uncompress the compressed file in trash and restore
	if err := utility.Uncompress(row.TrashPath, row.FromPath); err != nil {
		return err
	}
	
	// Delete unused entries in database
	if err := database.DeleteById(command.Id); err != nil {
		return err
	}

	return nil
}