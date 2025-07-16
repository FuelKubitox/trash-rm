package parser

import (
	"fmt"
	"strings"
)

/*
	Actions
	--------------------
	list - List all deleted files/folders
	delete - Delete a file or folder
	restore - Restore a file or folder
	empty - Empties the the trash
	help - Shows the help
*/

type Command struct {
	Action string
	Parameters []string
	Tags []string
	Target string
	Destination string
}

// Check all arguments and parse them. If everything is fine return true else false
func Parse(args []string) Command {
	command := Command{"", []string{}, []string{}, "", ""}

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
		default:
			fmt.Println("Wrong arguments. You can show the help with trm help")
	}

	return command
}

// List all objects in the trash
func listCommand(args []string, command Command) Command {
	if len(args) == 1 {
		command.Action = "list"
	} else if len(args) == 3 {
		if args[1] == "-t" {
			command.Action = "list"
			command.Parameters = []string{"t"}
			command.Tags = strings.Split(args[2], ",")
		} else {
			command.Action = "wrongArguments"
		}
	} else {
		command.Action = "wrongArguments"
	}
	return command
}

// Move a file or folder to the trash
func deleteCommand(args []string, command Command) Command {
	return command
}

// Restore an object in the trash
func restoreCommand(args []string, command Command) Command {
	return command
}

// Delete the objects in the trash and free space
func emptyCommand(args []string, command Command) Command {
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