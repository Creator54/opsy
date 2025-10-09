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
		m.logList.Title = m.getPathContext()
		logContent := m.logList.View()
		helpBar := m.renderLogsHelpBar()
		content = logContent + "\n\n" + helpBar
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
	
	var helpText string
	if m.showHelp {
		// Extended help
		helpText = "↑↓:Navigate | Enter:Open | g:Parent | H:Home | ?:Help | q:Quit"
	} else {
		// Short help
		helpText = "↑/k up · ↓/j down · enter open · g parent · H home · ? more · q quit"
	}
	return helpStyle.Render(helpText)
}

// renderLogsHelpBar renders help bar for logs mode
func (m model) renderLogsHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	var helpText string
	if m.showHelp {
		// Extended help
		helpText = "↑↓:Navigate | Enter:View | ?:Help | q:Back"
	} else {
		// Short help
		helpText = "↑/k up · ↓/j down · enter view · ? more · q back"
	}
	return helpStyle.Render(helpText)
}

// renderExecuteHelpBar renders help bar for execute mode
func (m model) renderExecuteHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	var helpText string
	if m.showHelp {
		// Extended help
		helpText = "↑↓/j/k:Navigate | Ctrl+u/d:Scroll | Enter:Run | e:Edit | s:Skip | l:SaveLog | ?:Help | q:Back"
	} else {
		// Short help
		helpText = "↑/k up · ↓/j down · ctrl+u/d scroll · enter run · e edit · s skip · l log · ? more · q back"
	}
	return helpStyle.Render(helpText)
}

// renderEditHelpBar renders help bar for edit mode
func (m model) renderEditHelpBar() string {
	helpStyle := statusBarStyle.Copy().
		Width(m.width).
		Foreground(colorFaint)
	
	var helpText string
	if m.showHelp {
		// Extended help
		helpText = "Enter:Save | Esc:Cancel | ?:Help"
	} else {
		// Short help
		helpText = "enter save · esc cancel · ? more"
	}
	return helpStyle.Render(helpText)
}
