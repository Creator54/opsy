package tui

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderExecutionContent generates the full content for the execution view
// Returns both the content and the line number where the current step starts
func (m model) renderExecutionContent() (content string, currentStepLine int) {
	if m.sop == nil {
		return "No SOP loaded", 0
	}

	var builder strings.Builder
	lineCount := 0

	// SOP Document Header
	sopTitle := m.sop.Title
	if sopTitle == "" {
		sopTitle = filepath.Base(m.sop.Path)
	}

	titleHeader := renderTitleHeader(sopTitle, m.width)
	builder.WriteString(titleHeader)
	lineCount += strings.Count(titleHeader, "\n")

	// Progress bar
	totalSteps := len(m.steps)
	completedSteps := 0
	for _, step := range m.steps {
		if step.Status == "success" {
			completedSteps++
		}
	}

	progressBar := renderProgressBar(completedSteps, totalSteps, m.width)
	builder.WriteString(progressBar)
	lineCount += strings.Count(progressBar, "\n")

	// Process each step with improved formatting
	for i, step := range m.steps {
		// Record the line number where current step starts
		if i == m.currentStep {
			currentStepLine = lineCount
		}

		isCurrent := i == m.currentStep

		// Step header
		stepHeader := renderStepHeader(i+1, step.Title, isCurrent)
		builder.WriteString(stepHeader + "\n")
		lineCount++

		// Status badge
		statusBadge := renderStatusBadge(step.Status, isCurrent, true)
		builder.WriteString(statusBadge + "\n\n")
		lineCount += 2

		// Description with better formatting
		if step.Description != "" {
			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("250")).
				PaddingLeft(4).
				Width(m.width - 8)

			wrappedDesc := wrapText(step.Description, m.width-12)
			descLines := strings.Count(wrappedDesc, "\n") + 1
			builder.WriteString(descStyle.Render(wrappedDesc) + "\n\n")
			lineCount += descLines + 1
		}

		// Command block
		if step.Command != "" {
			cmdBlock := renderCommandBlock(step.Command, m.width)
			builder.WriteString(cmdBlock)
			lineCount += strings.Count(cmdBlock, "\n")
		}

		// Output section
		if step.Output != "" {
			outputBlock := renderOutputBlock(step.Output, m.width, 8)
			builder.WriteString(outputBlock)
			lineCount += strings.Count(outputBlock, "\n")
		}

		// Error section
		if step.Error != "" {
			errorBlock := renderErrorBlock(step.Error, m.width, 5)
			builder.WriteString(errorBlock)
			lineCount += strings.Count(errorBlock, "\n")
		}

		// Step separator
		if i < len(m.steps)-1 {
			separator := renderStepSeparator(m.width)
			builder.WriteString(separator)
			lineCount += strings.Count(separator, "\n")
		}
	}

	return builder.String(), currentStepLine
}

// updateViewportContent updates the viewport with new content and scrolls to current step
func (m *model) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content, currentStepLine := m.renderExecutionContent()
	m.viewport.SetContent(content)

	// Only auto-scroll if manual scrolling is not active
	if !m.manualScrollActive {
		// Scroll to show the current step at the top of the viewport
		// This ensures the current step is always visible and easy to find
		targetOffset := currentStepLine
		if targetOffset < 0 {
			targetOffset = 0
		}
		// Ensure we don't scroll past the end of content
		maxOffset := m.viewport.TotalLineCount() - m.viewport.Height
		if maxOffset < 0 {
			maxOffset = 0
		}
		if targetOffset > maxOffset {
			targetOffset = maxOffset
		}

		m.viewport.SetYOffset(targetOffset)
	}
}
