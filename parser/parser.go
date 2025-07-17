package parser

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
)

/*
	Actions
	----------------------------------------------------------------------------
	list - List all deleted files/folders
	delete - Delete a file or folder
	restore - Restore a file or folder
	empty - Empties the the trash
	sync - Sync database with trash on filesystem and delete unnecassary entries
	help - Shows the help
*/
type Command struct {
	Action string
	Parameters []string
	Tags []string
	Id int
	Target string
	Destination string
}

// Contains the absolute path from the directory were we executed the program
var absolutePath string

// Check all arguments and parse them. Returns a Command object to execute the commands later
func Parse(argsFull []string) (Command, error) {
	var err error

	// Create an empty command object
	command := Command{"", []string{}, []string{}, 0, "", ""}

	// First get the absolute path from the dirctory were we executed the program.
	// We need that later to join the path with the path from the user.
	var path string
	path, err = filepath.Abs(filepath.Dir(argsFull[0]))
	if err != nil {
		return command, errors.New("cant get the absolute path from the directory you executed the program")
	}
	absolutePath = path

	// Filter the arguments and filter the first element (its usually the program name)
	args := argsFull[1:]

	// Choose and create command object with action
	switch args[0] {
		case "list":
			if command, err = listCommand(args, command); err != nil {
				return command, err
			}
		case "delete":
			if command, err = deleteCommand(args, command); err != nil {
				return command, err
			}
		case "restore":
			if command, err = restoreCommand(args, command); err != nil {
				return command, err
			}
		case "empty":
			if command, err = emptyCommand(args, command); err != nil {
				return command, err
			}
		case "help":
			if command, err = helpCommand(args, command); err != nil {
				return command, err
			}
		case "sync":
			if command, err = syncCommand(args, command); err != nil {
				return command, err
			}
		default:
			return command, errors.New("wrong arguments. You can show the help with trm help")
	}

	return command, nil
}

// List all objects in the trash
func listCommand(args []string, command Command) (Command, error) {
	if len(args) == 1 {
		// If only the list argument was passed
		command.Action = "list"
	} else if len(args) == 3 {
		// When you want to list objects in the trash by tags
		if args[1] == "-t" {
			command.Action = "list"
			command.Parameters = []string{"t"}
			command.Tags = strings.Split(args[2], ",")
		} else {
			// When the tag parameter doesnt match
			return command, errors.New("wrong arguments")
		}
	} else {
		// If the amount of arguments are wrong
		return command, errors.New("wrong arguments")
	}
	return command, nil
}

// Move a file or folder to the trash
func deleteCommand(args []string, command Command) (Command, error) {
	if len(args) == 2 {
		// Basic delete with only the file/folder name
		command.Action = "delete"
		command.Target = filepath.Join(absolutePath, args[1])
	} else if len(args) == 3 {
		// Delete with compression
		if args[1] == "-c" {
			command.Action = "delete"
			command.Parameters = []string{"c"}
			command.Target = filepath.Join(absolutePath, args[2])
		}
	} else if len(args) == 4 {
		// Delete with tags
		if args[1] == "-t" {
			command.Action = "delete"
			command.Parameters = []string{"t"}
			command.Tags = strings.Split(args[2], ",")
			command.Target = filepath.Join(absolutePath, args[3])
		} else {
			return command, errors.New("wrong arguments")
		}
	} else if len(args) == 5 {
		// Delete with tags and compression
		if args[1] == "-t" && args[3] == "-c" {
			command.Action = "delete"
			command.Parameters = []string{"t", "c"}
			command.Tags = strings.Split(args[2], ",")
			command.Target = filepath.Join(absolutePath, args[4])
		} else if args[1] == "-c" && args[2] == "-t" {
			command.Action = "delete"
			command.Parameters = []string{"t", "c"}
			command.Tags = strings.Split(args[3], ",")
			command.Target = filepath.Join(absolutePath, args[4])
		} else {
			return command, errors.New("wrong arguments")
		}
	} else {
		return command, errors.New("wrong arguments")
	}
	return command, nil
}

// Restore an object in the trash
func restoreCommand(args []string, command Command) (Command, error) {
	if len(args) == 2 {
		// Restore by id
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return command, errors.New("passed id is not a number")
		}
		command.Action = "restore"
		command.Id = id
	} else if len(args) == 3 {
		// Restore all objects in trash with tag
		if args[1] == "-t" {
			command.Action = "restore"
			command.Parameters = []string{"t"}
			command.Tags = strings.Split(args[2], ",")
		} else {
			return command, errors.New("wrong arguments")
		}
	} else if len(args) == 4 {
		// Restore by id with destination
		if args[2] == "-d" {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				return command, errors.New("passed id is not a number")
			}
			command.Action = "restore"
			command.Parameters = []string{"d"}
			command.Id = id
			command.Destination = filepath.Join(absolutePath, args[3])
		} else {
			return command, errors.New("wrongArguments")
		}
	} else {
		return command, errors.New("wrongArguments")
	}
	return command, nil
}

// Delete the objects in the trash and free space
func emptyCommand(args []string, command Command) (Command, error) {
	return command, nil
}

// Sync the database with the trash on the filesystem
func syncCommand(args []string, command Command) (Command, error) {
	return command, nil
}

// Show the help
func helpCommand(args []string, command Command) (Command, error) {
	if len(args) == 1 {
		command.Action = "help"
	} else {
		return command, errors.New("wrongArguments")
	}
	return command, nil
}