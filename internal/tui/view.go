package tui

import (
	"fmt"
)

// View renders the UI
func (m model) View() string {
	if m.quitting {
		return "Thanks for using Opsy! 👋\n"
	}

	var content string

	// Header section - shows "Opsy" and mode context (1 line with padding)
	headerText := "Opsy - " + m.getModeContext()
	header := headerStyle.Copy().Width(m.width).Render(headerText)

	// Main content section based on mode
	switch m.mode {
	case modeBrowse:
		// Update list title to show path context
		m.fileList.Title = m.getPathContext()
		listContent := m.fileList.View()
		helpBar := m.renderBrowseHelpBar()
		content = listContent + "\n\n" + helpBar
	case modeExecute:
		// Use viewport for execution view with help bar
		if m.viewportReady {
			viewportContent := m.viewport.View()
			// Add spacing and help bar for execute mode
			helpBar := m.renderExecuteHelpBar()
			content = viewportContent + "\n\n" + helpBar
		} else {
			content = "Loading..."
		}
	case modeLogs:
		if m.logViewReady {
			// Show log file content in viewport
			logViewContent := m.logViewPort.View()
			helpBar := m.renderLogsHelpBar()
			content = logViewContent + "\n\n" + helpBar
		} else {
			// Show log list
			m.logList.Title = m.getPathContext()
			logContent := m.logList.View()
			helpBar := m.renderLogsHelpBar()
			content = logContent + "\n\n" + helpBar
		}
	case modeEdit:
		// Show edit view with help bar
		editContent := m.renderEditView()
		helpBar := m.renderEditHelpBar()
		content = editContent + "\n\n" + helpBar
	}

	return header + "\n" + content
}

// renderEditView renders the edit mode view
func (m model) renderEditView() string {
	content := fmt.Sprintf("Editing Step %d/%d\n", m.currentStep+1, len(m.steps))
	content += m.textInput.View()
	return content
}

// renderBrowseHelpBar renders help bar for browse mode
func (m model) renderBrowseHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	// Short help only - consistent, concise text
	helpText := "↑↓ nav · ←/bs back · enter select · h home · l logs · q quit"
	return helpStyle.Render(helpText)
}

// renderLogsHelpBar renders help bar for logs mode
func (m model) renderLogsHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	// Short help only - consistent, concise text
	if m.logViewReady {
		// Help text when viewing a log file
		helpText := "↑↓ nav · q back"
		return helpStyle.Render(helpText)
	}
	
	// Help text when browsing log files
	helpText := "↑↓ nav · ←/bs back · enter select · q back"
	return helpStyle.Render(helpText)
}

// renderExecuteHelpBar renders help bar for execute mode
func (m model) renderExecuteHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	// Short help only - consistent, concise text
	helpText := "↑↓ nav · enter run · e edit · s skip · l logs · q back"
	return helpStyle.Render(helpText)
}

// renderEditHelpBar renders help bar for edit mode
func (m model) renderEditHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	// Short help only - consistent, concise text
	helpText := "enter save · esc cancel"
	return helpStyle.Render(helpText)
}
