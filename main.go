package main

import (
	"flag"
	"fmt"

	"mod-scavenger-file-manager/manual"
	// Updated import
)

var verbose bool

func main() {
	// Parse command-line flags
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	// modsDir := "./mods"

	// Define directories and lock file paths
	serverModsDir := "./manualmods/server"
	clientModsDir := "./manualmods/client"

	lockFilePath := "./mods.lock"

	// Update server mods
	if err := manual.UpdateMods(serverModsDir, lockFilePath, verbose); err != nil {
		fmt.Println("Error updating server mods:", err)
	}

	// Update client mods
	if err := manual.UpdateMods(clientModsDir, lockFilePath, verbose); err != nil {
		fmt.Println("Error updating client mods:", err)
	}
}
