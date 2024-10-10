package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Calculate the maximum length of mod names for proper alignment
// func getMaxModLength(modFiles []string, maxLimit int) int {
// 	maxLen := 0
// 	for _, modFile := range modFiles {
// 		modName := filepath.Base(modFile)
// 		if len(modName) > maxLen {
// 			maxLen = len(modName)
// 		}
// 	}
// 	if maxLen > maxLimit {
// 		return maxLimit
// 	}
// 	return maxLen
// }

// Truncate mod names if too long
func truncateModName(modName string, maxLength int) string {
	if len(modName) > maxLength {
		return modName[:maxLength-3] + "..."
	}
	return modName
}

// PrintTableStart prints the header and starts the table.
func PrintTableStart() {
	fmt.Println("+--------------------------------+-------------+------------+")
	fmt.Printf("| %-30s | %-11s | %-10s |\n", "Mod Name", "Environment", "Status")
	fmt.Println("+--------------------------------+-------------+------------+")
}

// PrintTableEnd closes the table.
func PrintTableEnd() {
	fmt.Println("+--------------------------------+-------------+------------+")
}

// PrintModTableEntry prints one line of the table for the given mod, environment, and status.
func PrintModTableEntry(modName, environment, status string) {
	// Define the environment styles using lipgloss for Bubbletea's color system
	envStyles := map[string]lipgloss.Style{
		"client":   lipgloss.NewStyle().Foreground(lipgloss.Color("42")),  // Green
		"server":   lipgloss.NewStyle().Foreground(lipgloss.Color("33")),  // Blue
		"agnostic": lipgloss.NewStyle().Foreground(lipgloss.Color("227")), // Yellow
	}

	// Define the mod name color styles based on status
	statusStyles := map[string]lipgloss.Style{
		"Added":    lipgloss.NewStyle().Foreground(lipgloss.Color("42")),  // Green for Added
		"Removed":  lipgloss.NewStyle().Foreground(lipgloss.Color("196")), // Red for Removed
		"Retained": lipgloss.NewStyle().Foreground(lipgloss.Color("15")),  // Default (white) for Retained
	}

	// Default style if environment not in map
	envStyle, exists := envStyles[environment]
	if !exists {
		envStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // Default to white
	}

	// Default style if status not in map
	modNameStyle, exists := statusStyles[status]
	if !exists {
		modNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // Default to white
	}

	// Print a single table row

	fmt.Printf("| %-30s | %-11s | %-10s |\n", modNameStyle.Render(fmt.Sprintf("%-30s", truncateModName(modName, 30))), envStyle.Render(fmt.Sprintf("%-11s", environment)), modNameStyle.Render(fmt.Sprintf("%-10s", status)))
}
