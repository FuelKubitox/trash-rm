package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

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

		// Check if the file in the trash exists and change the filename if yes
		destFile = checkIfDestExists(destFile)

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
		destBase := path.Base(command.Target)

		// Get index at leat occurence from point
		index := strings.LastIndex(destBase, ".")
	
		// Change suffix to .gz
		destBase = destBase[:index] + ".gz"

		// Define the destination path for the trash
		destFile := filepath.Join(trashDir, destBase)

		// Check if the file in the trash exists and change the filename if yes
		destFile = checkIfDestExists(destFile)
		
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
		i := strings.LastIndex(file, ".")
		dest = path.Join(dir, file[:i] + strconv.Itoa(count) + file[i:])
		_, err = os.Stat(dest)
	}
	return dest
}