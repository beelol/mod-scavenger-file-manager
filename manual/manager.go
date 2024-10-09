package manual

import (
	"fmt"
	"mod-scavenger-file-manager/lockfile"
	"mod-scavenger-file-manager/virtualization"
	"path/filepath"
)

// UpdateMods handles updating mods for the specified path (client or server)
func UpdateMods(modsDir, lockFilePath string, verbose bool) error {
	// Load the lock file
	lock, err := lockfile.LoadLockFile(lockFilePath)
	if err != nil {
		return err
	}

	// Track the mods that should be in the lock file after the update
	newLock := lockfile.LockFile{}

	// Get all mods in the specified directory
	modFiles, _ := filepath.Glob(filepath.Join(modsDir, "*.jar"))

	// Process mods, managing symlinks
	virtualization.ProcessSymlinks(lock, modFiles, modsDir, &newLock, verbose)

	// Save the updated lock file
	err = lockfile.SaveLockFile(lockFilePath, newLock)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %v", err)
	}

	return nil
}
