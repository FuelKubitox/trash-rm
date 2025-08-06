package commands

import (
	"database/sql"
	"fmt"

	"trash-rm/database"
	"trash-rm/parser"

	"github.com/jedib0t/go-pretty/v6/table"
)

func ListCommand(command parser.Command) error {
	var result *sql.Rows
	var trashList []database.TrashRow
	var err error
	if len(command.Parameters) == 0 {
		// If we want to list all trash objects from the database
		result, err = database.SelectAll()
		if err != nil {
			fmt.Println("Couldnt get data from db")
			return err
		}
	} else if len(command.Parameters) == 1 && command.Parameters[0] == "-t" && len(command.Tags) > 0 {
		// if we want to list trash objects filtered by tags from the database
		result, err = database.SelectWithTags(command.Tags)
		if err != nil {
			return err
		}
	}
	trashList, err = database.ParseDbData(result)
	if err != nil {
		return err
	}
	if err = showTrashList(trashList); err != nil {
		return err
	}
	return nil
}

// Show the results form the array struct TrashRow
func showTrashList(trashList []database.TrashRow) error {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Trash")
    t.AppendHeader(table.Row{"ID", "Basename", "From", "To", "Tags", "Deleted"})
	for _, trash := range trashList {
		t.AppendRow(table.Row{
			trash.Id,
			trash.Basename,
			trash.FromPath,
			trash.TrashPath,
			trash.Tags,
			trash.DeletedAt.Local(),
		})
	}
    fmt.Println(t.Render())
	return nil
}