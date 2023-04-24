package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Exists() will return true is a file exists
func Exists(path string) bool {
	// Get the file stats
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	// Check error
	if os.IsNotExist(err) {
		return false
	}

	return false
}

// Append() will append text to a file
func Append(path string, text string) error {
	// Open file
	f, err := os.OpenFile(path, os.O_APPEND, 0755)
	if err != nil {
		return fmt.Errorf("could not open file '%s'", path)
	}
	defer f.Close()

	// Write to file
	_, err = f.WriteString(text)
	if err != nil {
		return fmt.Errorf("could not write to file '%s'", path)
	}

	return nil
}

// Create() will create a file
func Create(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file '%s'", path)
	}
	f.Close()

	return nil
}

// ExeDir() returns the folder of the currently running executable
func ExeDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get path of executable with: %s", err)
	}

	return filepath.Dir(exe), nil
}

// ExePath() returns the full path of the executable
func ExePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get path of executable with: %s", err)
	}

	return exe, nil
}

// Copy() will copy a file from src to dst
//
// https://opensource.com/article/18/6/copying-files-go
func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// Open the src file
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// Create the dst file
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	// Copy content from src to dst
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
