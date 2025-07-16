package main

import (
	"fmt"
	"os"

	"trash-rm/parser"
)


var command parser.Command

func main() {
	command = parser.Parse(os.Args[1:])

	switch command.Action {
		case "list":
			fmt.Println("List of all files and folder in the trash")
		case "delete":
			fmt.Println("Delete a file or folder")
		case "restore":
			fmt.Println("Restore an object in the trash")
		case "empty":
			fmt.Println("Empty the trash and free space")
		case "help":
			fmt.Println("Show help")
		case "wrongArguments":
			fmt.Println("Wrong arguments! To show help use: trm help")
		default:
			fmt.Println("Unknow command!")
	}
}