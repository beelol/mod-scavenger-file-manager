package manual

import (
	"fmt"
	"mod-scavenger-file-manager/lockfile"
	"mod-scavenger-file-manager/ui"
	"mod-scavenger-file-manager/virtualization"
	"path/filepath"
)

// UpdateMods handles updating mods for the specified path (client or server)
func UpdateMods(modsDir, destDir, lockFilePath string, environment string, verbose bool) (lockfile.LockFile, error) {
	// Load the lock file
	lock, err := lockfile.LoadLockFile(lockFilePath)

	if err != nil {
		return lockfile.LockFile{}, err
	}

	// Track the mods that should be in the lock file after the update
	newLock := lockfile.LockFile{Mods: []lockfile.ModEntry{}}

	// Get all mods in the specified directory
	modFiles, _ := filepath.Glob(filepath.Join(modsDir, "*.jar"))

	// Display UI headers and progress
	// ui.DisplayHeader(fmt.Sprintf("Updating Mods for Environment: %s", environment))

	// ui.StartProgress(modFiles)

	// Process symlinks, ensuring that mods in the lockfile are still valid
	unlinkedModFiles := virtualization.ProcessSymlinks(lock, modFiles, modsDir, destDir, &newLock, verbose, environment)

	// Add new mods that are not already symlinked (unlinked mods)
	addNewModsToLockfile(unlinkedModFiles, &newLock, environment, modsDir, destDir, verbose)

	// fmt.Printf("Mods in newLock before saving: %+v\n", newLock.Mods)

	// ui.EndUI()
	return newLock, nil
}

// addNewModsToLockfile adds mods that are found in the directory but not in the lockfile
func addNewModsToLockfile(modFiles []string, newLock *lockfile.LockFile, environment, modsDir, destDir string, verbose bool) {
	for _, modFile := range modFiles {
		modName := filepath.Base(modFile) // Get the mod name

		// ui.UpdateProgress(modName)

		// Create a new ModEntry for each unlinked mod and add it to the new lockfile
		newMod := lockfile.ModEntry{
			Name:        modName,
			Version:     "1.0.0", // Set default version for now; this could be improved later
			Environment: environment,
			Source:      "local",
			FilePath:    modFile,
		}

		sourcePath := modFile
		destPath := filepath.Join(destDir, modName) // Adjust the destination as needed

		// Use the AddSymlink function to create the symlink
		if err := virtualization.AddSymlink(sourcePath, destPath, verbose); err != nil {
			fmt.Printf("Error creating symlink for %s: %v, not added to lockfile.\n", modName, err)
		} else {
			newLock.Mods = append(newLock.Mods, newMod)

			ui.PrintModTableEntry(modName, environment, "Added")
		}
	}
}
