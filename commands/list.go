package commands

import (
	"database/sql"
	"fmt"

	"trash-rm/database"
	"trash-rm/parser"
)

type TrashRow struct {
	Id int
	Basename string
	FromPath string
	TrashPath string
	CreatedAt string
	Tags []string
}

func ListCommand(command parser.Command) error {
	var result *sql.Rows
	var trashList []TrashRow
	var err error
	if len(command.Parameters) == 0 {
		// If we want to list all trash objects from the database
		result, err = listTrashAll()
		if err != nil {
			return err
		}
	} else if len(command.Parameters) == 1 && len(command.Tags) > 0 {
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
	selectAll := "SELECT trash_id, basename, from_path, trash_path, created_at" +
				" FROM trash_table RIGHT JOIN tags_table" +
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
	for result.Next() {
		trashRow := &TrashRow{}
		err := result.Scan(
			&trashRow.Id,
			&trashRow.Basename,
			&trashRow.FromPath,
			&trashRow.TrashPath,
			&trashRow.CreatedAt,
			&trashRow.Tags,
		)
		if err != nil {
			fmt.Println("Couldnt apply selected data to struct TrashRow")
			return trashList, err
		}
		trashList = append(trashList, *trashRow)
	}
	return trashList, nil
}

// Show the results form the array struct TrashRow
func showTrashList(trashList []TrashRow) error {
	for _, trash := range trashList {
		fmt.Printf("%+v", trash)
	}
	return nil
}