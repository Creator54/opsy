package tui

import (
	"fmt"
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

	// SOP Document Header with better styling
	sopTitle := m.sop.Title
	if sopTitle == "" {
		sopTitle = filepath.Base(m.sop.Path)
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width - 4).
		MarginBottom(1)

	builder.WriteString(titleStyle.Render(sopTitle) + "\n")
	lineCount++

	// Horizontal divider
	divider := strings.Repeat("─", m.width-4)
	dividerStyle := lipgloss.NewStyle().
		Foreground(colorBorder)
	builder.WriteString(dividerStyle.Render(divider) + "\n\n")
	lineCount += 2

	// Progress bar with visual indicator
	totalSteps := len(m.steps)
	completedSteps := 0
	for _, step := range m.steps {
		if step.Status == "success" {
			completedSteps++
		}
	}

	progressPercent := 0
	if totalSteps > 0 {
		progressPercent = (completedSteps * 100) / totalSteps
	}

	// Create visual progress bar
	barWidth := 40
	if m.width < 80 {
		barWidth = m.width - 40
		if barWidth < 10 {
			barWidth = 10
		}
	}
	filledWidth := (progressPercent * barWidth) / 100
	progressBar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

	progressStyle := lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true)

	progressLabel := fmt.Sprintf("Progress: %d/%d (%d%%)", completedSteps, totalSteps, progressPercent)
	builder.WriteString(progressStyle.Render(progressLabel) + "\n")
	builder.WriteString(progressBar + "\n\n")
	lineCount += 3

	// Process each step with improved formatting
	for i, step := range m.steps {
		// Record the line number where current step starts
		if i == m.currentStep {
			currentStepLine = lineCount
		}

		// Step container with border for current step
		isCurrent := i == m.currentStep

		// Step number and title on same line with better styling
		stepNumStyle := lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

		stepTitleStyle := lipgloss.NewStyle().
			Foreground(colorText).
			Bold(isCurrent)

		if isCurrent {
			stepTitleStyle = stepTitleStyle.Foreground(colorPrimary)
		}

		stepTitle := step.Title
		if stepTitle == "" {
			stepTitle = "Untitled Step"
		}

		// Add visual indicator for current step
		indicator := "  "
		if isCurrent {
			indicator = "▶ "
		}

		stepHeader := fmt.Sprintf("%s%s %s",
			indicator,
			stepNumStyle.Render(fmt.Sprintf("Step %d:", i+1)),
			stepTitleStyle.Render(stepTitle))

		builder.WriteString(stepHeader + "\n")
		lineCount++

		// Status badge with better styling
		var statusBadge string
		switch step.Status {
		case "success":
			badge := lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(colorSuccess).
				Padding(0, 1).
				Bold(true).
				Render("✓ DONE")
			statusBadge = "  " + badge
		case "error":
			badge := lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(colorError).
				Padding(0, 1).
				Bold(true).
				Render("✗ ERROR")
			statusBadge = "  " + badge
		case "skipped":
			badge := lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(colorWarning).
				Padding(0, 1).
				Bold(true).
				Render("⊘ SKIPPED")
			statusBadge = "  " + badge
		default:
			if isCurrent {
				badge := lipgloss.NewStyle().
					Foreground(lipgloss.Color("0")).
					Background(colorAccent).
					Padding(0, 1).
					Bold(true).
					Render("▶ CURRENT")
				statusBadge = "  " + badge
			} else {
				badge := lipgloss.NewStyle().
					Foreground(colorFaint).
					Background(lipgloss.Color("235")).
					Padding(0, 1).
					Render("○ PENDING")
				statusBadge = "  " + badge
			}
		}

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

		// Command block with improved styling
		if step.Command != "" {
			cmdLabelStyle := lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true).
				PaddingLeft(4)

			builder.WriteString(cmdLabelStyle.Render("Command:") + "\n")
			lineCount++

			commandBoxStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("117")).
				Background(lipgloss.Color("236")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(0, 1).
				MarginLeft(4).
				Width(m.width - 12)

			wrappedCmd := wrapText("$ "+step.Command, m.width-16)
			cmdLines := strings.Count(wrappedCmd, "\n") + 3 // +3 for border
			builder.WriteString(commandBoxStyle.Render(wrappedCmd) + "\n\n")
			lineCount += cmdLines + 1
		}

		// Output section with collapsible preview
		if step.Output != "" {
			outputLabelStyle := lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true).
				PaddingLeft(4)

			builder.WriteString(outputLabelStyle.Render("Output:") + "\n")
			lineCount++

			outputBoxStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Background(lipgloss.Color("234")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238")).
				Padding(0, 1).
				MarginLeft(4).
				Width(m.width - 12)

			truncatedOutput := truncateOutput(step.Output, 8)
			outputLines := strings.Count(truncatedOutput, "\n") + 3 // +3 for border
			builder.WriteString(outputBoxStyle.Render(truncatedOutput) + "\n\n")
			lineCount += outputLines + 1
		}

		// Error section with prominent styling
		if step.Error != "" {
			errorLabelStyle := lipgloss.NewStyle().
				Foreground(colorError).
				Bold(true).
				PaddingLeft(4)

			builder.WriteString(errorLabelStyle.Render("Error:") + "\n")
			lineCount++

			errorBoxStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("224")).
				Background(lipgloss.Color("52")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorError).
				Padding(0, 1).
				MarginLeft(4).
				Width(m.width - 12)

			truncatedError := truncateOutput(step.Error, 5)
			errorLines := strings.Count(truncatedError, "\n") + 3 // +3 for border
			builder.WriteString(errorBoxStyle.Render(truncatedError) + "\n\n")
			lineCount += errorLines + 1
		}

		// Step separator
		if i < len(m.steps)-1 {
			separator := lipgloss.NewStyle().
				Foreground(colorFaint).
				Render(strings.Repeat("─", m.width-4))
			builder.WriteString(separator + "\n\n")
			lineCount += 2
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

	// Scroll to show the current step
	// Position it roughly 1/4 from the top of the viewport for context
	targetOffset := currentStepLine - (m.viewport.Height / 4)
	if targetOffset < 0 {
		targetOffset = 0
	}

	m.viewport.SetYOffset(targetOffset)
}
