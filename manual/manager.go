package manual

import (
	"fmt"
	"mod-scavenger-file-manager/lockfile"
	"mod-scavenger-file-manager/ui"
	"mod-scavenger-file-manager/virtualization"
	"path/filepath"
)

// UpdateMods handles updating mods for the specified path (client or server)
func UpdateMods(modsDir, lockFilePath string, environment string, verbose bool) error {
	// Load the lock file
	lock, err := lockfile.LoadLockFile(lockFilePath)
	if err != nil {
		return err
	}

	// Track the mods that should be in the lock file after the update
	newLock := lockfile.LockFile{}

	// Get all mods in the specified directory
	modFiles, _ := filepath.Glob(filepath.Join(modsDir, "*.jar"))

	// Display UI headers and progress
	ui.DisplayHeader(fmt.Sprintf("Updating Mods for Environment: %s", environment))
	ui.StartProgress(modFiles)

	// Process symlinks, ensuring that mods in the lockfile are still valid
	unlinkedModFiles := virtualization.ProcessSymlinks(lock, modFiles, modsDir, &newLock, verbose, environment)

	// Add new mods that are not already symlinked (unlinked mods)
	addNewModsToLockfile(unlinkedModFiles, &newLock, environment, modsDir)

	// Save the updated lock file
	err = lockfile.SaveLockFile(lockFilePath, newLock)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %v", err)
	}

	ui.EndUI()
	return nil
}

// addNewModsToLockfile adds mods that are found in the directory but not in the lockfile
func addNewModsToLockfile(modFiles []string, newLock *lockfile.LockFile, environment, modsDir string) {
	for _, modFile := range modFiles {
		modName := filepath.Base(modFile) // Get the mod name

		ui.UpdateProgress(modName)

		// Create a new ModEntry for each unlinked mod and add it to the new lockfile
		newMod := lockfile.ModEntry{
			Name:        modName,
			Version:     "1.0.0", // Set default version for now; this could be improved later
			Environment: environment,
			Source:      "local",
			FilePath:    filepath.Join(modsDir, modName),
		}

		newLock.Mods = append(newLock.Mods, newMod)
	}
}
