package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"opsy/cmd"
	"opsy/internal/executor"
	"opsy/internal/logger"
	"opsy/internal/tui"
)

func main() {
	// Initialize executor
	executor := executor.NewExecutor()
	
	// Initialize logger
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger: ", err)
	}
	
	// Check if a command was provided
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			cmd.ListSOPs()
			return
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Usage: opsy [list]")
			os.Exit(1)
		}
	}
	
	// Default: launch TUI
	model := tui.NewModel(executor, logger)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program: ", err)
	}
}