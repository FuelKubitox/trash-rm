package main

import (
	"fmt"
	"os"

	"trash-rm/commands"
	"trash-rm/parser"
)

var command parser.Command
var err error

func main() {
	if command, err = parser.Parse(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch command.Action {
		case "list":
			fmt.Println("Listing all files and folder in trash")
		case "delete":
			fmt.Println("Start moving target(s) to trash...")
			if err := commands.DeleteCommand(command); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case "restore":
			fmt.Println("Restore an object in the trash")
		case "empty":
			fmt.Println("Empty the trash and free space")
		case "help":
			fmt.Println("Show help")
		case "wrongArguments":
			fmt.Println("Wrong arguments! To show help use: trm help")
		default:
			fmt.Println("Unknow action!")
			os.Exit(1)
	}
}