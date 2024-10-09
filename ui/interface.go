package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the UI, including headers and progress
type model struct {
	progress progress.Model
	files    []string
	index    int
	done     bool
	header   string
	details  []string
}

// Global model to track progress outside of async func
var uiModel *model

// Init is called when the UI starts
func (m model) Init() tea.Cmd {
	return tick()
}

// tick updates the UI at each step
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return t
	})
}

// Update handles state transitions in Bubbletea
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit
	}

	return m, nil
}

// View renders the header, progress bar, and details
func (m model) View() string {
	if m.done {
		return "All mods processed.\n"
	}

	file := m.files[m.index]
	progress := float64(m.index+1) / float64(len(m.files))
	header := fmt.Sprintf("%s\n", m.header)
	details := fmt.Sprintf("Processing: %s\n", file)
	return fmt.Sprintf("%s%s\n%s\n", header, details, m.progress.ViewAs(progress))
}

// DisplayHeader sets the header for the UI
func DisplayHeader(header string) {
	fmt.Printf("\n==== %s ====\n\n", header)
}

// StartProgress initializes the progress model
func StartProgress(files []string) {
	uiModel = &model{
		progress: progress.New(progress.WithDefaultGradient()),
		files:    files,
		header:   "Mod Update Progress",
		details:  []string{},
	}
}

// UpdateProgress updates the progress and logs the current file being added
func UpdateProgress(currentFile string) {
	uiModel.index++
	fmt.Printf("Adding %s\n", currentFile)

	if uiModel.index >= len(uiModel.files) {
		uiModel.done = true
	}
}

// EndUI ends the UI session once the operation is completed
func EndUI() {
	fmt.Println("\n==== Mod Update Complete ====\n")
}
