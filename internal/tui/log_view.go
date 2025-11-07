package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderLogExecutionView generates the full content for log viewing (similar to execute mode)
// Returns both the content and the line number where the current step starts
func (m model) renderLogExecutionView() (content string, currentStepLine int) {
	if m.logMetadata.Title == "" {
		return "No log loaded", 0
	}

	var builder strings.Builder
	lineCount := 0

	// Log title header with SOP name
	sopName := extractSOPName(m.logMetadata.SOPPath)
	title := fmt.Sprintf("%s (Logs)", sopName)
	titleHeader := renderTitleHeader(title, m.width)
	builder.WriteString(titleHeader)
	lineCount += strings.Count(titleHeader, "\n")

	// Metadata section
	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		PaddingLeft(4)

	if m.logMetadata.ExecutedBy != "" {
		builder.WriteString(metaStyle.Render(fmt.Sprintf("Executed by: %s", m.logMetadata.ExecutedBy)) + "\n")
		lineCount++
	}
	if m.logMetadata.StartedAt != "" {
		builder.WriteString(metaStyle.Render(fmt.Sprintf("Started: %s", m.logMetadata.StartedAt)) + "\n")
		lineCount++
	}
	if m.logMetadata.Status != "" {
		statusStyle := metaStyle.Copy()
		if strings.Contains(m.logMetadata.Status, "✅") || strings.Contains(m.logMetadata.Status, "Success") {
			statusStyle = statusStyle.Foreground(colorSuccess)
		} else if strings.Contains(m.logMetadata.Status, "❌") || strings.Contains(m.logMetadata.Status, "Failed") {
			statusStyle = statusStyle.Foreground(colorError)
		}
		builder.WriteString(statusStyle.Render(fmt.Sprintf("Status: %s", m.logMetadata.Status)) + "\n")
		lineCount++
	}
	builder.WriteString("\n")
	lineCount++

	// Progress indicator
	totalSteps := len(m.logSteps)
	completedSteps := 0
	for _, step := range m.logSteps {
		if strings.Contains(step.Status, "✅") || strings.Contains(step.Status, "Success") {
			completedSteps++
		}
	}

	progressBar := renderProgressBar(completedSteps, totalSteps, m.width)
	builder.WriteString(progressBar)
	lineCount += strings.Count(progressBar, "\n")

	// Render each step
	for i, step := range m.logSteps {
		// Record the line number where current step starts
		if i == m.currentLogStep {
			currentStepLine = lineCount
		}

		isCurrent := i == m.currentLogStep

		// Step header
		stepHeader := renderStepHeader(step.StepNumber, step.Title, isCurrent)
		builder.WriteString(stepHeader + "\n")
		lineCount++

		// Status badge (unified function handles both execute and log modes)
		statusBadge := renderStatusBadge(step.Status, isCurrent, false)
		builder.WriteString(statusBadge + "\n\n")
		lineCount += 2

		// Command block
		if step.Command != "" {
			cmdBlock := renderCommandBlock(step.Command, m.width)
			builder.WriteString(cmdBlock)
			lineCount += strings.Count(cmdBlock, "\n")
		}

		// Execution timestamp
		if step.ExecutedAt != "" {
			timeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				PaddingLeft(4)
			builder.WriteString(timeStyle.Render(fmt.Sprintf("Executed: %s", step.ExecutedAt)) + "\n\n")
			lineCount += 2
		}

		// Output section
		if step.HasOutput && step.Output != "" {
			outputBlock := renderOutputBlock(step.Output, m.width, 10)
			builder.WriteString(outputBlock)
			lineCount += strings.Count(outputBlock, "\n")
		}

		// Step separator
		if i < len(m.logSteps)-1 {
			separator := renderStepSeparator(m.width)
			builder.WriteString(separator)
			lineCount += strings.Count(separator, "\n")
		}
	}

	return builder.String(), currentStepLine
}

// updateLogViewportContent updates the log viewport with new content and scrolls to current step
func (m *model) updateLogViewportContent() {
	if !m.logViewReady {
		return
	}

	content, currentStepLine := m.renderLogExecutionView()
	m.logViewPort.SetContent(content)

	// Always scroll to show the current step at the top
	targetOffset := currentStepLine
	if targetOffset < 0 {
		targetOffset = 0
	}
	// Ensure we don't scroll past the end of content
	maxOffset := m.logViewPort.TotalLineCount() - m.logViewPort.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if targetOffset > maxOffset {
		targetOffset = maxOffset
	}

	m.logViewPort.SetYOffset(targetOffset)
}
