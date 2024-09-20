package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Verbose flag to control detailed output
var verbose bool

// UpdateMods reads the lock file, manages symlinks, and updates the lock file
func updateMods(serverModsDir, modsDir, lockFilePath string) {
	// Load the current lock file
	currentLock, err := loadLockFile(lockFilePath)

	if err != nil {
		if verbose {
			fmt.Printf("Lock file not found, creating one. Error: %v\n", err)
		}
		currentLock = make(map[string]bool) // Start fresh if lock file does not exist
	}

	// Track the mods that should be in the lock file after the update
	newLock := make(map[string]bool)

	// Iterate through the current lock to check for removals
	for modName := range currentLock {
		serverModPath := filepath.Join(serverModsDir, modName)
		modPath := filepath.Join(modsDir, modName)

		// Check if the mod still exists in the server-mods directory
		if _, err := os.Stat(serverModPath); os.IsNotExist(err) {
			// If it no longer exists in server-mods, remove its symlink from mods
			if info, err := os.Lstat(modPath); err == nil && (info.Mode()&os.ModeSymlink != 0) {
				if err := os.Remove(modPath); err != nil {
					fmt.Printf("Failed to remove symlink: %v\n", err)
				} else {
					if verbose {
						fmt.Printf("Removed outdated symlink: %s\n", modPath)
					}
				}
			}
		} else {
			// If it still exists in server-mods, keep it in the new lock
			newLock[modName] = true
		}
	}

	// Create symlinks for new or updated server mods
	serverModFiles, _ := filepath.Glob(filepath.Join(serverModsDir, "*.jar"))
	for _, serverMod := range serverModFiles {
		baseName := filepath.Base(serverMod)
		dest := filepath.Join(modsDir, baseName)

		// If the mod is not already in the new lock, it needs a symlink
		if _, found := newLock[baseName]; !found {
			if err := os.Symlink(serverMod, dest); err != nil {
				fmt.Printf("Failed to create symlink for %s: %v\n", serverMod, err)
			} else {
				if verbose {
					fmt.Printf("Symlink created for: %s\n", serverMod)
				}
				newLock[baseName] = true
			}
		}
	}

	// Save the updated lock file
	saveLockFile(lockFilePath, newLock)
}

// LoadLockFile loads the list of mods from the lock file
func loadLockFile(filePath string) (map[string]bool, error) {
	lock := make(map[string]bool)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return lock, nil // If file doesn't exist, return an empty lock
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lock[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lock, nil
}

// SaveLockFile saves the list of mods to the lock file
func saveLockFile(filePath string, lock map[string]bool) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for mod := range lock {
		_, err := writer.WriteString(mod + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func main() {
	// Parse command-line flags
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	// Define the directories and lock file path
	serverModsDir := "../server-mods"
	modsDir := "../mods"
	lockFilePath := "./server-mods.lock"

	// Update server mods
	updateMods(serverModsDir, modsDir, lockFilePath)
}
