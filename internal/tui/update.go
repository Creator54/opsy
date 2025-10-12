package tui

import (
	"fmt"
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
		// Refresh file list when returning to browse
		if msg.mode == modeBrowse {
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
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
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
	case "g": // Go to parent directory
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
	case "H": // Go to home (base) directory
		baseDir := config.DefaultBaseDirectory()
		cmds = append(cmds, func() tea.Msg {
			return browseToDirMsg{path: baseDir}
		})
		m.status = "Returned to base directory"
	case "q": // In browse mode, 'q' should quit the app
		m.quitting = true
		return m, tea.Quit
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
			// Update viewport content to reflect new current step
			m.updateViewportContent()
		} else {
			m.status = "Already at top"
		}
	case "down", "j": // Navigate down to next step
		if m.currentStep < len(m.steps)-1 {
			m.currentStep++
			m.status = fmt.Sprintf("Moved to step %d", m.currentStep+1)
			// Update viewport content to reflect new current step
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

		// Also let viewport handle other keys (pgup, pgdown, etc.)
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	return *m, tea.Batch(cmds...)
}

// handleLogKeys handles key events in log mode
func (m model) handleLogKeys(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q": // 'q' in log mode goes back to browse
		cmds = append(cmds, func() tea.Msg {
			return enterModeMsg{
				mode:   modeBrowse,
				status: "Returned to SOP browser",
			}
		})
	default:
		var cmd tea.Cmd
		m.logList, cmd = m.logList.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// handleEditKeys handles key events in edit mode
func (m model) handleEditKeys(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Save the edited command
		if m.currentStep < len(m.steps) {
			m.steps[m.currentStep].Command = m.textInput.Value()
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
		// Save execution log (implementation from handlers)
		cmds = append(cmds, m.saveExecutionLog())
	}
	return cmds
}
