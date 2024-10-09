package lockfile

import (
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
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(lock)
	if err != nil {
		return err
	}

	return nil
}
