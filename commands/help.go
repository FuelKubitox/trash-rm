package commands

import "fmt"

func HelpCommand() {
	fmt.Println("trm is a tool to soft delete files or folders. It compress the file/folder and moves it to a trash folder, so")
	fmt.Println("you can recover that file/folder later.")
	fmt.Printf("\n")
	fmt.Println("Followed arguments are supported")
	fmt.Println("----------------------------------------------------------------------------------------------")
	fmt.Println("trm delete [file/folder] - Moves the file or folder to the trash folder.")
	fmt.Println("trm delete -t [tags] [file/folder] - Same as above but you can add one or more tags. More tags are seperated with comma like 'tag1,tag2'.")
	fmt.Println("trm restore [id] - Restore the deleted file or folder from the trash folder. You need to pass the id as identifier.")
	fmt.Println("trm restore -t [tags] [id] - Same as above but only restore files and folders with the tags that was passed.")
	fmt.Println("trm list - List all deleted files and folders with informations like id.")
	fmt.Println("trm list -t [tags] - Same as above but filter by tags.")
	fmt.Println("trm empty - Delete all files and folders that are moved to the trash folders. It deletes every content in the trash folders.")
	fmt.Println("trm sync - Search for entries in the database that dont point to a deleted file or folder. It synchronizes the database with the trash folders.")
}