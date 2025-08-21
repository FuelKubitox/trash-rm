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
	// With no arguments do nothing
	if len(os.Args) <= 1 {
		fmt.Println("Missing arguments. For more information enter trm help")
		os.Exit(0)
	}

	// Parse passed arguments
	if command, err = parser.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Init database and create database object
	if err := database.InitDB(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Execute passed arguments with the command object
	switch command.Action {
		case "list":
			fmt.Println("Listing trash...")
			if err := commands.ListCommand(command); err != nil {
				fmt.Println(err)
			}
		case "delete":
			fmt.Println("Start moving file/folder to trash...")
			if err := commands.DeleteCommand(command); err != nil {
				fmt.Println(err)
			}
		case "restore":
			fmt.Println("Restore an object in the trash...")
			if err := commands.RestoreCommand(command); err != nil {
				fmt.Println(err)
			}
		case "empty":
			fmt.Println("Empty the trash and free space...")
			if err := commands.EmptyCommand(command); err != nil {
				fmt.Println(err)
			}
		case "sync":
			fmt.Println("Sync the database with trashes on the filesystem...")
			if err := commands.SyncCommand(command); err != nil {
				fmt.Println(err)
			}
		case "help":
			commands.HelpCommand()
		default:
			fmt.Println("Unknow action!")
	}

	defer database.Db.Close()
}