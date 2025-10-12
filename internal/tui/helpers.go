package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"

	"opsy/internal/config"
	"opsy/internal/parser"
)

// calculateContentHeight calculates the available height for content
// All modes now have custom help bars with consistent spacing
func calculateContentHeight(totalHeight int, isExecuteMode bool) int {
	// All modes use: header(1) + newline(1) + spacing(2) + helpbar(1) = 5 lines
	contentHeight := totalHeight - executeUIChromeHeight
	
	if contentHeight < minContentHeight {
		contentHeight = minContentHeight
	}
	return contentHeight
}

// wrapText wraps text to fit within specified width
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var wrapped strings.Builder
	lineLength := 0

	for _, word := range words {
		wordLength := len(word)

		// If word is longer than max width, break it
		if wordLength > maxWidth {
			if lineLength > 0 {
				wrapped.WriteString("\n")
				lineLength = 0
			}

			// Break long word into chunks
			for len(word) > maxWidth {
				wrapped.WriteString(word[:maxWidth])
				wrapped.WriteString("\n")
				word = word[maxWidth:]
			}

			if len(word) > 0 {
				wrapped.WriteString(word)
				lineLength = len(word)
			}
			continue
		}

		// Check if word fits on current line
		if lineLength+wordLength+1 <= maxWidth || lineLength == 0 {
			if lineLength > 0 {
				wrapped.WriteString(" ")
				lineLength++
			}
			wrapped.WriteString(word)
			lineLength += wordLength
		} else {
			// Word doesn't fit, start new line
			wrapped.WriteString("\n")
			wrapped.WriteString(word)
			lineLength = wordLength
		}
	}

	return wrapped.String()
}

// truncateOutput truncates output to specified number of lines
func truncateOutput(output string, maxLines int) string {
	lines := strings.Split(output, "\n")
	if len(lines) <= maxLines {
		return output
	}

	// Take first maxLines lines and add truncation indicator
	truncated := lines[:maxLines]
	truncated = append(truncated, fmt.Sprintf("... (%d more lines)", len(lines)-maxLines))
	return strings.Join(truncated, "\n")
}

// buildFileList builds a list of files and directories
func (m model) buildFileList(dir string) list.Model {
	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		// If error, go back to default directory
		dir = config.DefaultBaseDirectory()
		entries, err = os.ReadDir(dir)
		if err != nil {
			// Last resort - empty list
			entries = []os.DirEntry{}
		}
	}

	var items []list.Item

	// Add parent directory if not at the base directory
	defaultDir := config.DefaultBaseDirectory()
	if dir != defaultDir && dir != filepath.Dir(dir) { // Not at default or filesystem root
		parentDir := filepath.Dir(dir)
		items = append(items, item{
			title:    "../",
			desc:     "Parent directory",
			filePath: parentDir,
			isDir:    true,
		})
	}

	// Add directory and markdown file entries
	for _, entry := range entries {
		name := entry.Name()
		if name == ".." { // Skip the parent link we manually added
			continue
		}
		path := filepath.Join(dir, name)

		// Only include directories and markdown files
		if entry.IsDir() {
			items = append(items, item{
				title:    name + "/",
				desc:     "Directory",
				filePath: path,
				isDir:    true,
			})
		} else if strings.HasSuffix(strings.ToLower(name), ".md") {
			// Try to parse title from SOP
			title := name
			sop, err := parser.ParseSOP(path)
			if err == nil && sop.Title != "" {
				title = sop.Title
			}
			items = append(items, item{
				title:    name,
				desc:     title,
				filePath: path,
				isDir:    false,
			})
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 10)
	l.Title = "" // Remove title to keep it professional
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	return l
}

// buildLogList builds a list of log files and directories
func (m model) buildLogList(dir string) list.Model {
	// Read log directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		// If error, use log directory from config
		cfg := config.GetConfig()
		entries, err = os.ReadDir(cfg.LogDirectory)
		if err != nil {
			// Last resort - empty list
			entries = []os.DirEntry{}
		}
	}

	var items []list.Item

	// Add parent directory if not at the root log directory
	cfg := config.GetConfig()
	if dir != cfg.LogDirectory && dir != filepath.Dir(dir) { // Not at log root or filesystem root
		parentDir := filepath.Dir(dir)
		items = append(items, item{
			title:    "../",
			desc:     "Parent directory",
			filePath: parentDir,
			isDir:    true,
		})
	}

	// Separate directories and files to sort them properly
	var dirs []struct {
		entry os.DirEntry
		info  os.FileInfo
	}
	var files []struct {
		entry os.DirEntry
		info  os.FileInfo
	}

	// Get file info for sorting
	for _, entry := range entries {
		name := entry.Name()
		if name == ".." { // Skip the parent link we manually added
			continue
		}

		info, err := entry.Info()
		if err != nil {
			// Skip if we can't get file info
			continue
		}

		if entry.IsDir() {
			dirs = append(dirs, struct {
				entry os.DirEntry
				info  os.FileInfo
			}{entry: entry, info: info})
		} else if strings.HasSuffix(strings.ToLower(name), ".log.md") {
			files = append(files, struct {
				entry os.DirEntry
				info  os.FileInfo
			}{entry: entry, info: info})
		}
	}

	// Sort directories by modification time (newest first)
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].info.ModTime().After(dirs[j].info.ModTime())
	})

	// Sort files by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].info.ModTime().After(files[j].info.ModTime())
	})

	// Add directories first (sorted)
	for _, dirEntry := range dirs {
		name := dirEntry.entry.Name()
		
		// Better description for SOP directories
		desc := "SOP directory"
		
		path := filepath.Join(dir, name)
		items = append(items, item{
			title:    name + "/",
			desc:     desc,
			filePath: path,
			isDir:    true,
		})
	}

	// Add files (sorted)
	for _, fileEntry := range files {
		name := fileEntry.entry.Name()
		
		// Extract date and timestamp from filename for better description
		desc := "Execution log"
		// Try to extract date and timestamp from filename like "sop-name_DD-MM-YYYY_HH-MM-SS.log.md"
		parts := strings.Split(strings.TrimSuffix(name, ".log.md"), "_")
		if len(parts) >= 3 {
			// Format should be: sop-name_date_timestamp
			if len(parts) >= 2 {
				date := parts[len(parts)-2] // Second to last part should be the date
				timestamp := parts[len(parts)-1] // Last part should be the timestamp
				if len(date) >= 10 && date[2] == '-' && date[5] == '-' &&
				   len(timestamp) >= 8 && timestamp[2] == '-' && timestamp[5] == '-' {
					// Looks like a date DD-MM-YYYY and timestamp HH-MM-SS
					desc = fmt.Sprintf("Execution: %s %s", date, timestamp)
				}
			}
		}
		
		path := filepath.Join(dir, name)
		items = append(items, item{
			title:    name,
			desc:     desc,
			filePath: path,
			isDir:    false,
		})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 10)
	l.Title = "" // Remove title to keep it professional
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	return l
}

// getModeContext returns a string describing the current mode
func (m model) getModeContext() string {
	switch m.mode {
	case modeBrowse:
		return "Browser"
	case modeExecute:
		if m.sop != nil {
			return "Execution"
		}
		return "Execution"
	case modeLogs:
		return "Logs"
	case modeEdit:
		return "Edit"
	default:
		return ""
	}
}

// getPathContext returns context information for the current path
func (m model) getPathContext() string {
	switch m.mode {
	case modeBrowse:
		return m.currentPath
	case modeExecute:
		if m.sop != nil {
			return m.sop.Title
		}
		return ""
	case modeLogs:
		return "Execution logs directory"
	case modeEdit:
		return fmt.Sprintf("Step %d/%d", m.currentStep+1, len(m.steps))
	default:
		return ""
	}
}

// getHelpText returns help text for the current mode
func (m model) getHelpText() string {
	switch m.mode {
	case modeBrowse:
		return "↑↓ nav · ←/bs back · enter select · h home · l logs · q quit"
	case modeExecute:
		return "↑↓ nav · enter run · e edit · s skip · l logs · q back"
	case modeLogs:
		return "↑↓ nav · ←/bs back · enter select · q back"
	case modeEdit:
		return "enter save · esc cancel"
	default:
		return "q quit"
	}
}
