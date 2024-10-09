package lockfile

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ModEntry represents a single mod in the lockfile
type ModEntry struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"` // "client", "server", or "agnostic"
	Source      string `yaml:"source"`      // "local" or "remote"`
	FilePath    string `yaml:"file_path"`
}

// LockFile represents the structure of the lockfile
type LockFile struct {
	Mods []ModEntry `yaml:"mods"`
}

// LoadLockFile loads the list of mods from the YAML lock file
func LoadLockFile(filePath string) (LockFile, error) {
	var lock LockFile

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return lock, nil // If file doesn't exist, return an empty lockfile
		}
		return lock, err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&lock)
	if err != nil {
		return lock, err
	}

	return lock, nil
}

// SaveLockFile saves the list of mods to the YAML lock file
func SaveLockFile(filePath string, lock LockFile) error {
	// Create the file, overwrite if it exists
	file, err := os.Create(filePath)

	if err != nil {
		fmt.Printf("Error creating lockfile: %v\n", err)
		return err
	}

	defer file.Close()

	// Create a buffered writer to ensure everything is written before the file is closed
	writer := bufio.NewWriter(file)

	// Initialize the YAML encoder with the writer
	encoder := yaml.NewEncoder(writer)

	// Try encoding the lockfile into the writer
	err = encoder.Encode(lock)
	if err != nil {
		fmt.Printf("Error encoding lockfile to YAML: %v\n", err)
		return err
	}

	// Ensure everything is flushed from the writer to the file
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Error flushing lockfile to disk: %v\n", err)
		return err
	}

	// Check for final errors on file close
	err = file.Close()
	if err != nil {
		fmt.Printf("Error closing lockfile: %v\n", err)
		return err
	}

	fmt.Println("Lockfile saved successfully.")
	return nil
}
