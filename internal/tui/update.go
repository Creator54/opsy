package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"opsy/internal/config"
	"opsy/internal/parser"
)

// Update handles all state updates
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		
		// Calculate content height for lists
		listHeight := calculateContentHeight(m.height, false)
		m.fileList.SetSize(m.width, listHeight)
		m.fileList.SetShowHelp(false)
		m.fileList.SetShowPagination(false)
		m.logList.SetSize(m.width, listHeight)
		m.logList.SetShowHelp(false)
		m.logList.SetShowPagination(false)

		// Initialize or update viewport for execute mode
		if m.mode == modeExecute {
			viewportHeight := calculateContentHeight(m.height, true)
			if !m.viewportReady {
				m.viewport = viewport.New(m.width, viewportHeight)
				m.viewport.YPosition = 0
				m.viewportReady = true
				m.updateViewportContent()
			} else {
				m.viewport.Width = m.width
				m.viewport.Height = viewportHeight
				m.updateViewportContent()
			}
		}

	case browseToDirMsg:
		m.currentPath = msg.path
		m.fileList = m.buildFileList(m.currentPath)
		m.status = fmt.Sprintf("Directory changed to: %s", filepath.Base(msg.path))
		// Ensure the list is properly sized and configured
		if m.width > 0 && m.height > 0 {
			listHeight := calculateContentHeight(m.height, false)
			m.fileList.SetSize(m.width, listHeight)
			m.fileList.SetShowHelp(false)
			m.fileList.SetShowPagination(false)
		}

	case enterModeMsg:
		m.mode = msg.mode
		if msg.status != "" {
			m.status = msg.status
		}
		if msg.sop != nil {
			m.sop = msg.sop
		}
		if msg.steps != nil {
			m.steps = msg.steps
		}
		// Initialize viewport when entering execute mode
		if msg.mode == modeExecute {
			viewportHeight := calculateContentHeight(m.height, true)
			m.viewport = viewport.New(m.width, viewportHeight)
			m.viewport.YPosition = 0
			m.viewportReady = true
			m.updateViewportContent()
		}
		// Initialize and refresh log list when entering logs mode
		if msg.mode == modeLogs {
			cfg := config.GetConfig()
			// Store SOP path for returning to browse mode
			m.sopPath = m.currentPath
			// Store previous mode for proper return
			m.previousMode = msg.from
			// Determine which logs to show based on context
			logsPath := cfg.LogDirectory
			if msg.path != "" {
				// If we have a context path, try to find corresponding logs
				currentFolder := filepath.Base(msg.path)
				if currentFolder != "." && currentFolder != "/" {
					// Check if there's a logs subdirectory for this SOP folder
					contextLogsPath := filepath.Join(cfg.LogDirectory, currentFolder)
					if _, err := os.Stat(contextLogsPath); err == nil {
						logsPath = contextLogsPath
					}
				}
			}
			m.logList = m.buildLogList(logsPath)
			m.currentPath = logsPath // Update current path to logs path
			// Ensure proper sizing
			if m.width > 0 && m.height > 0 {
				listHeight := calculateContentHeight(m.height, false)
				m.logList.SetSize(m.width, listHeight)
				m.logList.SetShowHelp(false)
				m.logList.SetShowPagination(false)
			}
		}
		// Refresh file list when returning to browse
		if msg.mode == modeBrowse {
			// Restore SOP path if we were in logs mode
			if m.sopPath != "" {
				m.currentPath = m.sopPath
				m.sopPath = "" // Clear the stored SOP path
			}
			m.fileList = m.buildFileList(m.currentPath)
			m.viewportReady = false
			// Ensure proper sizing
			if m.width > 0 && m.height > 0 {
				listHeight := calculateContentHeight(m.height, false)
				m.fileList.SetSize(m.width, listHeight)
				m.fileList.SetShowHelp(false)
				m.fileList.SetShowPagination(false)
			}
		}

	case tea.KeyMsg:
		// Global key bindings
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		// Mode-specific key bindings
		switch m.mode {
		case modeBrowse:
			return m.handleBrowseKeys(msg, cmds)
		case modeExecute:
			return (&m).handleExecuteKeys(msg, cmds)
		case modeLogs:
			return m.handleLogKeys(msg, cmds)
		case modeEdit:
			return m.handleEditKeys(msg, cmds)
		}
	}

	// Update active list component
	switch m.mode {
	case modeBrowse:
		// Already handled above for Enter key, but handle other messages
		// Only update if we didn't handle Enter specially above
		if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
			m.fileList, cmd = m.fileList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case modeLogs:
		m.logList, cmd = m.logList.Update(msg)
		cmds = append(cmds, cmd)
	case modeExecute:
		// Viewport is updated in handleExecuteKeys
	}

	return m, tea.Batch(cmds...)
}

// handleBrowseKeys handles key events in browse mode
func (m model) handleBrowseKeys(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "h": // Go to home (base) directory
		baseDir := config.DefaultBaseDirectory()
		cmds = append(cmds, func() tea.Msg {
			return browseToDirMsg{path: baseDir}
		})
		m.status = "Returned to base directory"
	case "l": // Go to logs directory
		cmds = append(cmds, func() tea.Msg {
			return enterModeMsg{
				mode:   modeLogs,
				status: "Entered logs browser",
				path:   m.currentPath, // Pass current path for logs filtering
				from:   m.mode,        // Pass current mode as previous mode
			}
		})
		m.status = "Entering logs browser"
	case "q": // In browse mode, 'q' should quit the app
		m.quitting = true
		return m, tea.Quit
	case "backspace": // Go back to parent directory
		// Navigate to parent directory when backspace is pressed
		parentDir := filepath.Dir(m.currentPath)
		baseDir := config.DefaultBaseDirectory()

		if m.currentPath != baseDir && parentDir != m.currentPath {
			cmds = append(cmds, func() tea.Msg {
				return browseToDirMsg{path: parentDir}
			})
			m.status = "Moved to parent directory"
		} else {
			m.status = "Already at base directory"
		}
	case "enter": // Handle Enter to select item
		// Get the selected item from the list
		if selectedItem, ok := m.fileList.SelectedItem().(item); ok {
			// Check if this is a parent directory reference
			if selectedItem.title == "../" {
				// Handle parent directory case
				parentDir := filepath.Dir(m.currentPath)
				baseDir := config.DefaultBaseDirectory()

				if m.currentPath != baseDir && parentDir != m.currentPath {
					cmds = append(cmds, func() tea.Msg {
						return browseToDirMsg{path: parentDir}
					})
					m.status = "Moved to parent directory"
				} else {
					m.status = "Already at base directory"
				}
			} else if selectedItem.isDir {
				// Enter directory
				cmds = append(cmds, func() tea.Msg {
					return browseToDirMsg{path: selectedItem.filePath}
				})
				m.status = fmt.Sprintf("Entered: %s", selectedItem.title)
			} else {
				// Load SOP file
				sop, err := parser.ParseSOP(selectedItem.filePath)
				if err != nil {
					m.status = fmt.Sprintf("Error loading SOP: %v", err)
				} else {
					steps := make([]SOPStep, len(sop.Steps))
					for i, step := range sop.Steps {
						steps[i] = SOPStep{
							ID:          step.ID,
							Title:       step.Title,
							Description: step.Description,
							Command:     step.Command,
							Status:      "pending",
						}
					}
					cmds = append(cmds, func() tea.Msg {
						return enterModeMsg{
							mode:   modeExecute,
							status: fmt.Sprintf("Loaded SOP: %s", sop.Title),
							sop:    sop,
							steps:  steps,
						}
					})
				}
			}
		}
	default:
		// Let the list component handle navigation keys (up/down/etc.)
		var cmd tea.Cmd
		m.fileList, cmd = m.fileList.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	return m, tea.Batch(cmds...)
}

// handleExecuteKeys handles key events in execute mode
func (m *model) handleExecuteKeys(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q": // 'q' in execute mode goes back to browse
		cmds = append(cmds, func() tea.Msg {
			return enterModeMsg{
				mode:   modeBrowse,
				status: "Returned to SOP browser",
			}
		})
	case "up", "k": // Navigate up to previous step
		if m.currentStep > 0 {
			m.currentStep--
			m.status = fmt.Sprintf("Moved to step %d", m.currentStep+1)
			// Update viewport content to reflect new current step and scroll to it
			m.updateViewportContent()
		} else {
			m.status = "Already at top"
		}
	case "down", "j": // Navigate down to next step
		if m.currentStep < len(m.steps)-1 {
			m.currentStep++
			m.status = fmt.Sprintf("Moved to step %d", m.currentStep+1)
			// Update viewport content to reflect new current step and scroll to it
			m.updateViewportContent()
		} else {
			m.status = "Already at last step"
		}
	case "ctrl+u": // Scroll view up by half page
		m.viewport.HalfViewUp()
		m.status = "Scrolled up"
	case "ctrl+d": // Scroll view down by half page
		m.viewport.HalfViewDown()
		m.status = "Scrolled down"
	default:
		// Handle execute mode commands (enter, e, s, l)
		executeCmds := m.handleExecuteCommands(msg)
		cmds = append(cmds, executeCmds...)

		// Also let viewport handle other keys
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	return *m, tea.Batch(cmds...)
}

// handleLogKeys handles key events in log mode
func (m model) handleLogKeys(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q": // 'q' in log mode goes back to previous mode
		if m.logViewReady {
			// Exit log view mode and return to log list
			m.logViewReady = false
			m.logViewContent = ""
			m.logViewPath = ""
			m.status = "Returned to log list"
		} else {
			// Return to previous mode
			var returnMode string
			var status string
			
			// Determine which mode to return to based on where we came from
			switch m.previousMode {
			case modeExecute:
				returnMode = modeExecute
				status = "Returned to SOP execution"
			default:
				returnMode = modeBrowse
				status = "Returned to SOP browser"
			}
			
			cmds = append(cmds, func() tea.Msg {
				return enterModeMsg{
					mode:   returnMode,
					status: status,
				}
			})
		}
	case "h": // Go to home (base) directory from logs mode
		cmds = append(cmds, func() tea.Msg {
			return enterModeMsg{
				mode:   modeBrowse,
				status: "Returned to base directory",
			}
		})
		m.status = "Returned to base directory"
	case "backspace": // Go back to parent directory in logs mode
		// Navigate to parent directory when backspace is pressed
		parentDir := filepath.Dir(m.currentPath)
		cfg := config.GetConfig()
		logsDir := cfg.LogDirectory

		if m.currentPath != logsDir && parentDir != m.currentPath {
			// Build log list for parent directory
			m.logList = m.buildLogList(parentDir)
			m.currentPath = parentDir
			m.status = "Moved to parent directory"
			// Ensure proper sizing
			if m.width > 0 && m.height > 0 {
				listHeight := calculateContentHeight(m.height, false)
				m.logList.SetSize(m.width, listHeight)
				m.logList.SetShowHelp(false)
				m.logList.SetShowPagination(false)
			}
		} else {
			m.status = "Already at logs root directory"
		}
	case "enter": // Handle Enter to select item in log mode
		// Only process Enter if not already viewing a log file
		if !m.logViewReady {
			// Get the selected item from the list
			if selectedItem, ok := m.logList.SelectedItem().(item); ok {
				// Check if this is a parent directory reference
				if selectedItem.title == "../" {
					// Handle parent directory case
					cfg := config.GetConfig()
					parentDir := filepath.Dir(m.currentPath)
					
					// Only go back if not at the log root directory
					if m.currentPath != cfg.LogDirectory && parentDir != m.currentPath {
						// Build log list for parent directory
						m.logList = m.buildLogList(parentDir)
						m.currentPath = parentDir
						m.status = "Moved to parent directory"
						// Ensure proper sizing
						if m.width > 0 && m.height > 0 {
							listHeight := calculateContentHeight(m.height, false)
							m.logList.SetSize(m.width, listHeight)
							m.logList.SetShowHelp(false)
							m.logList.SetShowPagination(false)
						}
					} else {
						m.status = "Already at log root directory"
					}
				} else if selectedItem.isDir {
					// Enter directory
					m.logList = m.buildLogList(selectedItem.filePath)
					m.currentPath = selectedItem.filePath
					m.status = fmt.Sprintf("Entered: %s", selectedItem.title)
					// Ensure proper sizing
					if m.width > 0 && m.height > 0 {
						listHeight := calculateContentHeight(m.height, false)
						m.logList.SetSize(m.width, listHeight)
						m.logList.SetShowHelp(false)
						m.logList.SetShowPagination(false)
					}
				} else {
					// View log file content
					content, err := os.ReadFile(selectedItem.filePath)
					if err != nil {
						m.status = fmt.Sprintf("Error reading log file: %v", err)
					} else {
						// Store the log content and path
						m.logViewContent = string(content)
						m.logViewPath = selectedItem.filePath
						
						// Switch to a simple log viewing mode (using a viewport)
						viewportHeight := calculateContentHeight(m.height, true)
						m.logViewPort = viewport.New(m.width, viewportHeight)
						m.logViewPort.YPosition = 0
						m.logViewPort.SetContent(m.logViewContent)
						m.logViewReady = true
						m.status = fmt.Sprintf("Viewing log: %s", selectedItem.title)
					}
				}
			}
		}
	case "up", "k": // Scroll up
		if m.logViewReady {
			m.logViewPort.LineUp(1)
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "down", "j": // Scroll down
		if m.logViewReady {
			m.logViewPort.LineDown(1)
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "pgup": // Page up
		if m.logViewReady {
			m.logViewPort.PageUp()
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "pgdown", " ": // Page down or space
		if m.logViewReady {
			m.logViewPort.PageDown()
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "home": // Go to top
		if m.logViewReady {
			m.logViewPort.GotoTop()
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "ctrl+u": // Scroll view up by half page
		if m.logViewReady {
			m.logViewPort.HalfViewUp()
			m.status = "Scrolled up"
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "ctrl+d": // Scroll view down by half page
		if m.logViewReady {
			m.logViewPort.HalfViewDown()
			m.status = "Scrolled down"
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case "end": // Go to bottom
		if m.logViewReady {
			m.logViewPort.GotoBottom()
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	default:
		// Handle scrolling when viewing log content
		if m.logViewReady {
			var cmd tea.Cmd
			m.logViewPort, cmd = m.logViewPort.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.logList, cmd = m.logList.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

// handleEditKeys handles key events in edit mode
func (m model) handleEditKeys(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Save the edited command to both the temporary steps and original SOP
		if m.currentStep < len(m.steps) && m.currentStep < len(m.sop.Steps) {
			editedCommand := m.textInput.Value()
			m.steps[m.currentStep].Command = editedCommand
			// Also update the original SOP step so execution uses the edited command
			m.sop.Steps[m.currentStep].Command = editedCommand
			cmds = append(cmds, func() tea.Msg {
				return enterModeMsg{
					mode:   modeExecute,
					status: "Command updated",
				}
			})
			// Update viewport to show the change
			m.updateViewportContent()
		}
	case "esc":
		// Cancel editing
		cmds = append(cmds, func() tea.Msg {
			return enterModeMsg{
				mode:   modeExecute,
				status: "Edit cancelled",
			}
		})
	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// handleExecuteCommands handles command execution keys
func (m *model) handleExecuteCommands(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "enter", " ":
		// Run current step (no auto-advance)
		if m.currentStep < len(m.steps) {
			step := m.sop.Steps[m.currentStep]
			result, err := m.executor.ExecuteStep(step)
			if err != nil {
				m.steps[m.currentStep].Status = statusError
				m.steps[m.currentStep].Error = err.Error()
				m.status = fmt.Sprintf("Error executing step: %v", err)
			} else {
				m.steps[m.currentStep].Status = result.Status
				m.steps[m.currentStep].Output = result.Output
				if result.Status == statusSuccess {
					m.status = "Step executed successfully"
				} else {
					m.status = fmt.Sprintf("Step execution %s", result.Status)
				}
			}
			// Update viewport content to show execution results
			m.updateViewportContent()
			
			// Save incremental log after each command execution
			cmds = append(cmds, m.saveExecutionLog())
		}
	case "e":
		// Edit command
		if m.currentStep < len(m.steps) {
			m.textInput.SetValue(m.steps[m.currentStep].Command)
			cmds = append(cmds, func() tea.Msg {
				return enterModeMsg{
					mode:   modeEdit,
					status: m.status,
				}
			})
			cmds = append(cmds, textinput.Blink)
		}
	case "s":
		// Skip current step (no auto-advance)
		if m.currentStep < len(m.steps) {
			m.steps[m.currentStep].Status = statusSkipped
			m.status = "Step skipped"
			// Update viewport content to show skip
			m.updateViewportContent()
		}
	case "l":
		// Go to logs browser
		cmds = append(cmds, func() tea.Msg {
			return enterModeMsg{
				mode:   modeLogs,
				status: "Entered logs browser",
				path:   filepath.Dir(m.sop.Path), // Pass SOP directory path for logs filtering
				from:   m.mode,                  // Pass current mode as previous mode
			}
		})
	}
	return cmds
}
