package virtualization

import (
	"fmt"
	"mod-scavenger-file-manager/lockfile"
	"mod-scavenger-file-manager/ui"
	"os"
	"path/filepath"
)

// ProcessSymlinks manages the symlinking process for mods and returns unlinked mods, respecting the environment
func ProcessSymlinks(lock lockfile.LockFile, sourceModFiles []string, sourceDirectory, destinationDirectory string, newLock *lockfile.LockFile, verbose bool, environment string) []string {
	// Create a map of mods that are already symlinked
	symlinkedMods := map[string]bool{}

	// Iterate over the lockfile mods and check if they still exist in the source directory, matching the environment
	for _, mod := range lock.Mods {
		// Skip mods that don't match the environment
		if mod.Environment != environment {
			continue
		}

		sourceModPath := filepath.Join(sourceDirectory, filepath.Base(mod.FilePath))
		destModPath := filepath.Join(destinationDirectory, filepath.Base(mod.FilePath))

		// If the mod exists in the source directory, ensure symlink is created and keep it in the lockfile
		if modExists(sourceModFiles, mod.FilePath) {
			// newLock.Mods = append(newLock.Mods, mod)
			newLock.Mods = append(newLock.Mods, mod)

			symlinkedMods[mod.FilePath] = true

			// Ensure the symlink exists; recreate if necessary
			if _, err := os.Lstat(destModPath); os.IsNotExist(err) {
				err := AddSymlink(sourceModPath, destModPath, verbose)
				if err != nil {
					fmt.Printf("Error recreating symlink for %s: %v\n", sourceModPath, err)

					// ui.PrintModTableEntry(mod.Name, environment, "Failed Recreation")
				} else {
					if verbose {
						fmt.Printf("Recreated missing symlink for %s\n", sourceModPath)
					}

					ui.PrintModTableEntry(mod.Name, environment, "Recreated")
				}
			}

			ui.PrintModTableEntry(mod.Name, environment, "Retained")

		} else {
			// If the mod doesn't exist in the source directory, remove its symlink and exclude it from the lockfile
			removeSymlink(destModPath, verbose)

			ui.PrintModTableEntry(mod.Name, environment, "Removed")
		}
	}

	// Return the mods that are not yet symlinked (i.e., not in the lockfile)
	return filterUnlinkedMods(sourceModFiles, symlinkedMods)
}

// filterUnlinkedMods returns mods that are in the modFiles but not in symlinkedMods
func filterUnlinkedMods(modFiles []string, symlinkedMods map[string]bool) []string {
	var unlinkedMods []string
	for _, modFile := range modFiles {
		if !symlinkedMods[modFile] {
			unlinkedMods = append(unlinkedMods, modFile)
		}
	}
	return unlinkedMods
}

// modExists checks if the mod exists in the modFiles list
func modExists(modFiles []string, modPath string) bool {
	for _, modFile := range modFiles {
		if filepath.Base(modFile) == filepath.Base(modPath) {
			return true
		}
	}
	return false
}

// removeSymlink removes the symlink for the specified mod
func removeSymlink(modPath string, verbose bool) {
	if info, err := os.Lstat(modPath); err == nil && (info.Mode()&os.ModeSymlink != 0) {
		if err := os.Remove(modPath); err != nil {
			fmt.Printf("Failed to remove symlink: %v\n", err)
		} else if verbose {
			fmt.Printf("Removed outdated symlink: %s\n", modPath)
		}
	}
}

// AddSymlink creates a symbolic link from the source file to the destination
func AddSymlink(sourcePath, destPath string, verbose bool) error {
	if verbose {
		fmt.Printf("Attempting to create symlink from %s to %s\n", sourcePath, destPath)
	}

	// Check if the symlink already exists
	if _, err := os.Lstat(destPath); os.IsNotExist(err) {
		// Create the symlink if it doesn't exist
		err := os.Symlink(sourcePath, destPath)
		if err != nil {
			if verbose {
				fmt.Printf("Failed to create symlink for %s: %v\n", sourcePath, err)
			}
			return err
		} else {
			if verbose {
				fmt.Printf("Symlink created for: %s\n", sourcePath)
			}
		}
	} else if verbose {
		fmt.Printf("Symlink already exists for: %s\n", sourcePath)
	}
	return nil
}
