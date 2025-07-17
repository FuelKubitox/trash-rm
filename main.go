package main

import (
	"log"
	"os"

	"trash-rm/commands"
	"trash-rm/parser"
)

var command parser.Command
var err error

func main() {
	if command, err = parser.Parse(os.Args); err != nil {
		log.Fatal(err)
	}

	switch command.Action {
		case "list":
			log.Println("Listing all files and folder in trash")
		case "delete":
			log.Println("Start deleting")
			if err := commands.DeleteCommand(command); err != nil {
				log.Fatal(err)
			}
		case "restore":
			log.Println("Restore an object in the trash")
		case "empty":
			log.Println("Empty the trash and free space")
		case "help":
			log.Println("Show help")
		case "wrongArguments":
			log.Println("Wrong arguments! To show help use: trm help")
		default:
			log.Fatal("Unknow action!")
	}
}