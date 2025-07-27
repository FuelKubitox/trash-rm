package commands

import (
	"database/sql"
	"fmt"
	"time"

	"trash-rm/database"
	"trash-rm/parser"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TrashRow struct {
	Id int
	Basename string
	FromPath string
	TrashPath string
	CreatedAt time.Time
	Tags string
}

func ListCommand(command parser.Command) error {
	var result *sql.Rows
	var trashList []TrashRow
	var err error
	if len(command.Parameters) == 0 {
		// If we want to list all trash objects from the database
		result, err = listTrashAll()
		if err != nil {
			fmt.Println("Couldnt get data from db")
			return err
		}
	} else if len(command.Parameters) == 1 && command.Parameters[0] == "-t" && len(command.Tags) > 0 {
		// if we want to list trash objects filtered by tags from the database
		result, err = listTrashTags(command)
		if err != nil {
			return err
		}
	}
	trashList, err = parseDbData(result)
	if err != nil {
		return err
	}
	if err = showTrashList(trashList); err != nil {
		return err
	}
	return nil
}

// List all objects in trash without any filter
func listTrashAll() (*sql.Rows, error) {
	selectAll := "SELECT trash_id, basename, from_path, trash_path, created_at, GROUP_CONCAT(tags_table.tagname, ',') AS tags" +
				" FROM trash_table LEFT JOIN tags_table" +
				" ON trash_table.trash_id = tags_table.trash_tag_id"
	result, err := database.Db.Query(selectAll)
	if err != nil {
		fmt.Println("Couldnt select all trash objects from the database")
		return result, err
	}
	return result, nil
}

// List objects filtered by tags
func listTrashTags(command parser.Command) (*sql.Rows, error) {
	var result *sql.Rows
	return result, nil
}

// Take the query result and parse it to a TrashRow array struct
func parseDbData(result *sql.Rows) ([]TrashRow, error) {
	var trashList []TrashRow
	var row TrashRow
	var tags sql.NullString
	for result.Next() {
		err := result.Scan(
			&row.Id,
			&row.Basename,
			&row.FromPath,
			&row.TrashPath,
			&row.CreatedAt,
			&tags,
		)
		row.Tags = tags.String
		if err != nil {
			fmt.Println("Couldnt scan result from db")
			return trashList, err
		}
		trashList = append(trashList, row)
	}
	return trashList, nil
}

// Show the results form the array struct TrashRow
func showTrashList(trashList []TrashRow) error {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Trash")
    t.AppendHeader(table.Row{"ID", "Basename", "From", "To", "Tags", "Created"})
	for _, trash := range trashList {
		t.AppendRow(table.Row{
			trash.Id,
			trash.Basename,
			trash.FromPath,
			trash.TrashPath,
			trash.Tags,
			trash.CreatedAt.Local(),
		})
	}
    fmt.Println(t.Render())
	return nil
}