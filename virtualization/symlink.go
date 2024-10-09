package virtualization

import (
	"fmt"
	"mod-scavenger-file-manager/lockfile"
	"os"
	"path/filepath"
)

// ProcessSymlinks manages the symlinking process for mods
func ProcessSymlinks(lock lockfile.LockFile, modFiles []string, modsDir string, newLock *lockfile.LockFile, verbose bool) {
	for _, mod := range lock.Mods {
		modPath := filepath.Join(modsDir, filepath.Base(mod.FilePath))

		// Check if the mod still exists in the mod directory
		if !modExists(modFiles, mod.FilePath) {
			removeSymlink(modPath, verbose)
		} else {
			// If it still exists, add it to the new lockfile
			newLock.Mods = append(newLock.Mods, mod)
		}
	}
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
