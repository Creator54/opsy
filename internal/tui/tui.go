package tui

// This file serves as the main entry point for the TUI package.
// The implementation is split across multiple files for better maintainability:
//
// - model.go: Core model struct and initialization
// - styles.go: All styling constants and lipgloss styles
// - messages.go: Custom message types for state changes
// - update.go: Update logic and event handling
// - view.go: View rendering
// - execute_view.go: Execution view rendering with viewport
// - handlers.go: Command handlers (save log, etc.)
// - helpers.go: Utility functions (text wrapping, file listing, etc.)

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the TUI application
func Run(executor ExecutorInterface, logger LoggerInterface) error {
	p := tea.NewProgram(
		NewModel(executor, logger),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
