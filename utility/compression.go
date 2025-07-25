package utility

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Compress the file to the target location
// source: file to compress with absolute path to file
// dest: absolute path to the file in the trash with the suffix .gz at the end
func CompressFile(source string, dest string) error {
	// Open file
	sfile, err := os.Open(source)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt compress the file to trash")
	}
	defer sfile.Close()

	// Read bytes to buffer
	read := bufio.NewReader(sfile)

	// Read all bytes
	data, err := io.ReadAll(read)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt read all bytes out of the buffer from the file during compress")
	}

	// Create compressed file
	dfile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt create compressed file during compress")
	}

	// Open gzip writer
	write, err := gzip.NewWriterLevel(dfile, gzip.BestSpeed)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt open the gip writer during compress")
	}

	// Write compressed data to file
	write.Write(data)

	// Close the file writer
	write.Close()

	return nil
}

// Uncompress a file from trash
// source: file to uncompress with suffix .gz at the end
// dest: absolute path where the file will be uncompressed to
func UncompressFile(source string, dest string) error {
	// Open compressed file
	sfile, err := os.Open(source)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt open compressed file during uncompress")
	}
	defer sfile.Close()

	// Create reader for file
	reader := bufio.NewReader(sfile)
	
	// Read compressed data
	data, err := gzip.NewReader(reader)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt read compressed data from file during uncompress")
	}
	
	// Read all chunk data to buffer
	buffer, err := io.ReadAll(data)
	if err != nil {
		fmt.Println(err)
		return errors.New("coudlnt read all buffer data during uncompress")
	}

	// Create new file to write uncompressed data
	dfile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err)
		return errors.New("couldnt create file for uncompressed data")
	}
	defer dfile.Close()

	// Create writer to write uncompressed data to file
	writer := bufio.NewWriter(dfile)

	// Write uncompressed to file
	writer.Write(buffer)

	return nil
}

// Compress a directory to the target location recursively
// sourceDir: directory to compress with absolute path
// destinationFile: absolute path to the file in the trash with the suffix .gz at the end
func CompressDir(source string, dest string) error {
	// Count all files
    files, err := listFiles(source)
    if err != nil {
        return err
    }
    total := len(files)
    if total == 0 {
        return fmt.Errorf("now files in directory found: %s", source)
    }

    // Create target file during compress
    dfile, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer dfile.Close()

    // Create gzip writer
    gwriter, err := gzip.NewWriterLevel(dfile, gzip.BestSpeed)
    if err != nil {
		fmt.Println(err)
		return errors.New("couldnt create gzip writer during compress")
	}
	defer gwriter.Close()

    // Create tar writer
    twriter := tar.NewWriter(gwriter)
    defer twriter.Close()

    // Walk recursivly through the directory
    for i, file := range files {
		relPath, err := filepath.Rel(source, file)
        if err != nil {
            return err
        }

        if err := addFileToTar(twriter, file, relPath); err != nil {
            return err
        }

        // Fortschritt ausgeben
        printProgress(i+1, total, relPath)
	}

	return nil
}

// Add file to tar archive
func addFileToTar(tw *tar.Writer, path, relPath string) error {
    fi, err := os.Stat(path)
    if err != nil {
        return err
    }

    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()

    header, err := tar.FileInfoHeader(fi, "")
    if err != nil {
        return err
    }

    header.Name = filepath.ToSlash(relPath)

    if err := tw.WriteHeader(header); err != nil {
        return err
    }

    _, err = io.Copy(tw, f)
    return err
}

// Uncompress a file from trash
// source: file to uncompress with suffix .gz at the end
// dest: absolute path where the file will be uncompressed to
func UncompressDir(source string, dest string) error {
    // Open compressed file
    file, err := os.Open(source)
    if err != nil {
        return err
    }
    defer file.Close()

    // Create gzip reader
    greader, err := gzip.NewReader(file)
    if err != nil {
        return err
    }
    defer greader.Close()

    // Create tar reader
    treader := tar.NewReader(greader)

	// Total files
	totalFiles := 0

	entries, err := countTarEntries(source)
    if err != nil {
        return err
    }

	fmt.Printf("Restore %d entries...\n\n", entries)

    // Loop through the tar archive
    for {
        header, err := treader.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        // Create target path
        targetPath := filepath.Join(dest, header.Name)

        switch header.Typeflag {
        case tar.TypeDir:
            // Create directory
            if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
                return err
            }
        case tar.TypeReg:
            // Secure directory
            if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
                return err
            }

            // Create file
            outFile, err := os.Create(targetPath)
            if err != nil {
                return err
            }

            // Copy data
            if _, err := io.Copy(outFile, treader); err != nil {
                outFile.Close()
                return err
            }
            outFile.Close()

            // Set chmod
            if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
                return err
            }
        default:
            fmt.Printf("Skip unknown type: %v in %s\n", header.Typeflag, header.Name)
        }

		totalFiles++
		printProgress(totalFiles, entries, header.Name)
    }

    return nil
}

// Count all files recursively
func listFiles(root string) ([]string, error) {
    var files []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.Mode().IsRegular() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

// Print progress
func printProgress(current, total int, name string) {
    percent := float64(current) / float64(total) * 100
    fmt.Printf("[%3.0f%%] %s\n", percent, name)
}

// Count entries in archive during uncompress
func countTarEntries(archive string) (int, error) {
    f, err := os.Open(archive)
    if err != nil {
        return 0, err
    }
    defer f.Close()

    greader, err := gzip.NewReader(f)
    if err != nil {
        return 0, err
    }
    defer greader.Close()

    tarReader := tar.NewReader(greader)

    count := 0
    for {
        _, err := tarReader.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return 0, err
        }
        count++
    }

    return count, nil
}