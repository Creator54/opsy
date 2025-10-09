package cmd

// This package will contain the command implementations for the CLI

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"opsy/internal/config"
	"opsy/internal/parser"
)

// ListSOPs lists all available SOPs in the base directory
func ListSOPs() {
	cfg := config.GetConfig()
	
	// Check if base directory exists
	if _, err := os.Stat(cfg.BaseDirectory); os.IsNotExist(err) {
		fmt.Printf("Base directory does not exist: %s\n", cfg.BaseDirectory)
		fmt.Println("Please create the directory or update your configuration.")
		return
	}
	
	// Walk through the directory to find all .md files
	err := filepath.Walk(cfg.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Only process markdown files
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			// Try to parse the SOP to get its title
			sop, parseErr := parser.ParseSOP(path)
			if parseErr != nil {
				fmt.Printf("  [ERROR] %s: could not parse (%v)\n", path, parseErr)
				return nil // Continue with other files
			}
			
			fmt.Printf("  %s: %s\n", path, sop.Title)
		}
		
		return nil
	})
	
	if err != nil {
		log.Printf("Error walking through directory: %v", err)
	}
}