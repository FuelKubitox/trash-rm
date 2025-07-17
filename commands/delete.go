package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"trash-rm/parser"
	"trash-rm/utility"
)

const trash string = ".trash"

func DeleteCommand(command parser.Command) error {
	if len(command.Parameters) == 0 && command.Target != "" {
		if err := basicDelete(command); err != nil {
			return err
		}
	} else {
		return errors.New("some mismatch in the command object")
	}
	return nil
}

func basicDelete(command parser.Command) error {
	// First check if the given target exists
	var targetInfo os.FileInfo
	var err error
	if targetInfo, err = os.Stat(command.Target); os.IsNotExist(err) {
		return errors.New("file or folder that was passed doesnt exist")
	}
	
	// Get the absolute path were the target is on the filesystem
	dir := path.Dir(command.Target)
	
	// Create .trash dir if not exist
	trashDir := filepath.Join(dir, trash)
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		os.Mkdir(trashDir, os.ModePerm)
	}

	// Get only the target if a path was passed
	target := path.Base(command.Target)
	
	// Define the destination path for the trash
	destination := filepath.Join(trashDir, target)
	fmt.Println(destination)
	
	if targetInfo.IsDir() {
		// Move the target if its a directory
		if err = utility.MoveDirectory(make(chan int, 10), command.Target, destination); err != nil {
			return err
		}
	} else {
		// Move the target if its a file
		if err = utility.MoveFile(command.Target, destination); err != nil {
			return err
		}
	}
	return nil
}