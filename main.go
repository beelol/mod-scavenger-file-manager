package main

import (
	"flag"
	"fmt"

	"mod-scavenger-file-manager/lockfile"
	"mod-scavenger-file-manager/manual"
	"mod-scavenger-file-manager/ui"
)

var verbose bool

// mergeLockFiles merges mods from multiple lockfiles into a single lockfile
func mergeLockFiles(baseLock, newLock lockfile.LockFile) lockfile.LockFile {
	// Create a map to avoid duplicates based on mod name
	modMap := make(map[string]lockfile.ModEntry)

	// Add all existing mods from baseLock to the map
	for _, mod := range baseLock.Mods {
		modMap[mod.Name] = mod
	}

	// Add or update mods from newLock
	for _, mod := range newLock.Mods {
		modMap[mod.Name] = mod
	}

	// Convert the map back into a slice
	mergedMods := make([]lockfile.ModEntry, 0, len(modMap))
	for _, mod := range modMap {
		mergedMods = append(mergedMods, mod)
	}

	// Return the merged lockfile
	return lockfile.LockFile{Mods: mergedMods}
}

func main() {
	// Parse command-line flags
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	// modsDir := "./mods"

	// Define directories and lock file paths
	serverModsDir := "./manualmods/server"
	clientModsDir := "./manualmods/client"
	agnosticModsDir := "./manualmods"
	destinationModsDir := "./mods"

	lockFilePath := "./mods.lock"

	virtualLock := lockfile.LockFile{}

	ui.PrintTableStart()

	clientLock, err := manual.UpdateMods(clientModsDir, destinationModsDir, lockFilePath, "client", verbose)
	if err != nil {
		fmt.Printf("Error updating client mods: %v\n", err)
		return
	}

	virtualLock = mergeLockFiles(virtualLock, clientLock)

	serverLock, err := manual.UpdateMods(serverModsDir, destinationModsDir, lockFilePath, "server", verbose)
	if err != nil {
		fmt.Printf("Error updating server mods: %v\n", err)
		return
	}

	virtualLock = mergeLockFiles(virtualLock, serverLock)

	agnosticLock, err := manual.UpdateMods(agnosticModsDir, destinationModsDir, lockFilePath, "agnostic", verbose)
	if err != nil {
		fmt.Printf("Error updating agnostic mods: %v\n", err)
		return
	}

	ui.PrintTableEnd()

	virtualLock = mergeLockFiles(virtualLock, agnosticLock)

	err = lockfile.SaveLockFile(lockFilePath, virtualLock)
	if err != nil {
		fmt.Printf("Error saving lockfile: %v\n", err)
	} else {
		fmt.Println("Lockfile saved successfully.")
	}

}
