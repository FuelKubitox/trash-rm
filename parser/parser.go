package parser

import (
	"fmt"
	"strconv"
	"strings"
)

/*
	Actions
	--------------------
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

// Check all arguments and parse them. If everything is fine return true else false
func Parse(args []string) Command {
	command := Command{"", []string{}, []string{}, 0, "", ""}

	switch args[0] {
		case "list":
			command = listCommand(args, command)
		case "delete":
			command = deleteCommand(args, command)
		case "restore":
			command = restoreCommand(args, command)
		case "empty":
			command = emptyCommand(args, command)
		case "help":
			command = helpCommand(args, command)
		case "sync":
			command = syncCommand(args, command)
		default:
			fmt.Println("Wrong arguments. You can show the help with trm help")
	}

	return command
}

// List all objects in the trash
func listCommand(args []string, command Command) Command {
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
			command.Action = "wrongArguments"
		}
	} else {
		// If the amount of arguments are wrong
		command.Action = "wrongArguments"
	}
	return command
}

// Move a file or folder to the trash
func deleteCommand(args []string, command Command) Command {
	if len(args) == 2 {
		// Basic delete with only the file/folder name
		command.Action = "delete"
		command.Target = args[1]
	} else if len(args) == 3 {
		// Delete with compression
		if args[1] == "-c" {
			command.Action = "delete"
			command.Parameters = []string{"c"}
			command.Target = args[2]
		}
	} else if len(args) == 4 {
		// Delete with tags
		if args[1] == "-t" {
			command.Action = "delete"
			command.Parameters = []string{"t"}
			command.Tags = strings.Split(args[2], ",")
			command.Target = args[3]
		} else {
			command.Action = "wrongArguments"
		}
	} else if len(args) == 5 {
		// Delete with tags and compression
		if args[1] == "-t" && args[3] == "-c" {
			command.Action = "delete"
			command.Parameters = []string{"t", "c"}
			command.Tags = strings.Split(args[2], ",")
			command.Target = args[4]
		} else if args[1] == "-c" && args[2] == "-t" {
			command.Action = "delete"
			command.Parameters = []string{"t", "c"}
			command.Tags = strings.Split(args[3], ",")
			command.Target = args[4]
		} else {
			command.Action = "wrongArguments"
		}
	} else {
		command.Action = "wrongArguments"
	}
	return command
}

// Restore an object in the trash
func restoreCommand(args []string, command Command) Command {
	if len(args) == 2 {
		// Restore by id
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Passed id is not a number")
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
			command.Action = "wrongArguments"
		}
	} else if len(args) == 4 {
		// Restore by id with destination
		if args[2] == "-d" {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Passed id is not a number")
			}
			command.Action = "restore"
			command.Parameters = []string{"d"}
			command.Id = id
			command.Destination = args[3]
		} else {
			command.Action = "wrongArguments"
		}
	} else {
		command.Action = "wrongArguments"
	}
	return command
}

// Delete the objects in the trash and free space
func emptyCommand(args []string, command Command) Command {
	return command
}

// Sync the database with the trash on the filesystem
func syncCommand(args []string, command Command) Command {
	return command
}

// Show the help
func helpCommand(args []string, command Command) Command {
	if len(args) == 1 {
		command.Action = "help"
	} else {
		command.Action = "wrongArguments"
	}
	return command
}