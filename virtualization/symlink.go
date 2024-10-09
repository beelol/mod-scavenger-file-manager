package virtualization

import (
	"fmt"
	"mod-scavenger-file-manager/lockfile"
	"os"
	"path/filepath"
)

// ProcessSymlinks manages the symlinking process for mods and returns unlinked mods, respecting the environment
func ProcessSymlinks(lock lockfile.LockFile, modFiles []string, modsDir string, newLock *lockfile.LockFile, verbose bool, environment string) []string {
	// Create a map of mods that are already symlinked
	symlinkedMods := map[string]bool{}

	// Iterate over the lockfile mods and check if they still exist in the mod directory, matching the environment
	for _, mod := range lock.Mods {
		// Skip mods that don't match the environment
		if mod.Environment != environment {
			continue
		}

		modPath := filepath.Join(modsDir, filepath.Base(mod.FilePath))

		// If the mod exists in the mod directory, keep it in the new lockfile and mark as symlinked
		if modExists(modFiles, mod.FilePath) {
			newLock.Mods = append(newLock.Mods, mod)
			symlinkedMods[mod.FilePath] = true
		} else {
			// Otherwise, remove its symlink
			removeSymlink(modPath, verbose)
		}
	}

	// Return the mods that are not yet symlinked (i.e., not in the lockfile)
	return filterUnlinkedMods(modFiles, symlinkedMods)
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
