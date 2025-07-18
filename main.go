package main

import (
	"fmt"
	"os"

	"trash-rm/commands"
	"trash-rm/database"
	"trash-rm/parser"
)

// Command object contains everything to execute the program with its arguments
var command parser.Command
var err error

func main() {
	// Parse passed arguments
	if command, err = parser.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Init database
	if err := database.InitDB(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Execute passed arguments with the command object
	switch command.Action {
		case "list":
			fmt.Println("Listing all files and folder in trash")
		case "delete":
			fmt.Println("Start moving target(s) to trash...")
			if err := commands.DeleteCommand(command); err != nil {
				fmt.Println(err)
			}
		case "restore":
			fmt.Println("Restore an object in the trash")
		case "empty":
			fmt.Println("Empty the trash and free space")
		case "help":
			fmt.Println("Show help")
		default:
			fmt.Println("Unknow action!")
	}
}