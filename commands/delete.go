package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"trash-rm/database"
	"trash-rm/parser"
	"trash-rm/utility"
)

const trash string = ".trash"

func DeleteCommand(command parser.Command) error {
	if len(command.Parameters) == 0 && command.Target != "" {
		// Without any parameters
		if err := delete(command, []string{}); err != nil {
			return err
		}
	} else if len(command.Parameters) == 1 && command.Target != "" {
		// With parameters for tags
		if command.Parameters[0] != "-t" {
			return errors.New("unknown parameter")
		} else if len(command.Tags) == 0 {
			return errors.New("no tags are passed")
		} else {
			if err := delete(command, command.Tags); err != nil {
				return err
			}
		}
	} else {
		return errors.New("some mismatch in the command object")
	}
	return nil
}

func delete(command parser.Command, tags []string) error {
	// First check if the given target exists
	var targetInfo os.FileInfo
	var err error
	if targetInfo, err = os.Stat(command.Target); os.IsNotExist(err) {
		fmt.Println(err)
		return errors.New("file or folder that was passed doesnt exist")
	}
	
	// Get the absolute path were the target is on the filesystem
	dir := path.Dir(command.Target)
	
	// Create .trash dir if not exist
	trashDir := filepath.Join(dir, trash)
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		os.Mkdir(trashDir, os.ModePerm)
	}
	
	if targetInfo.IsDir() {
		baseDir := path.Base(command.Target)
		destFile := path.Join(trashDir, baseDir + ".gz")

		// Check if the file in the trash already exists and change the filename if yes
		destFile = checkIfDestExists(destFile)

		// Create the database entry
		if err := database.Insert(baseDir, command.Target, destFile, tags); err != nil {
			return err
		}

		// Compress the target if its a directory
		if err = utility.CompressDir(command.Target, destFile); err != nil {
			return err
		}

		// Remove the directory
		if err = os.RemoveAll(command.Target); err != nil {
			return err
		}
	} else {
		// Get only the target if a path was passed
		baseFile := path.Base(command.Target)

		// Get index at leat occurence from point
		index := strings.LastIndex(baseFile, ".")
	
		// Change suffix to .gz
		baseFile = baseFile[:index] + ".gz"

		// Define the destination path for the trash
		destFile := filepath.Join(trashDir, baseFile)

		// Check if the file in the trash already exists and change the filename if yes
		destFile = checkIfDestExists(destFile)
		
		// Create the database entry
		if err := database.Insert(baseFile, command.Target, destFile, tags); err != nil {
			return err
		}

		// Compress the target if its a file
		if err = utility.CompressFile(command.Target, destFile); err != nil {
			return err
		}

		// Delete source file
		if err := os.Remove(command.Target); err != nil {
			fmt.Println(err)
			return errors.New("couldnt remove source file after compression")
		}
	}

	return nil
}

// Check if the destination file exists, if yes change the filename and add a number
func checkIfDestExists(dest string) string {
	var err error
	// File count
	count := 1
	// Get file info
	_, err = os.Stat(dest)
	// If file info returns an error, then add the file count to the filename and loop until the file doesnt exist
	// So we always have a new file in the trash and dont overwrite an existing one
	for !os.IsNotExist(err) {
		dir := path.Dir(dest)
		file := path.Base(dest)
		i := strings.LastIndex(file, "_")
		if i == 0 {
			i = strings.LastIndex(file, ".")
		}
		dest = path.Join(dir, file[:i] + "_" + strconv.Itoa(count) + ".gz")
		
		_, err = os.Stat(dest)
	}
	return dest
}