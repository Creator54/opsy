package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderStepHeader renders a step header with number and title
func renderStepHeader(stepNum int, title string, isCurrent bool) string {
	stepNumStyle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	stepTitleStyle := lipgloss.NewStyle().
		Foreground(colorText).
		Bold(isCurrent)

	if isCurrent {
		stepTitleStyle = stepTitleStyle.Foreground(colorPrimary)
	}

	if title == "" {
		title = "Untitled Step"
	}

	// Add visual indicator for current step
	indicator := "  "
	if isCurrent {
		indicator = "▶ "
	}

	return fmt.Sprintf("%s%s %s",
		indicator,
		stepNumStyle.Render(fmt.Sprintf("Step %d:", stepNum)),
		stepTitleStyle.Render(title))
}

// renderCommandBlock renders a command in a styled box
func renderCommandBlock(command string, width int) string {
	if command == "" {
		return ""
	}

	var builder strings.Builder

	cmdLabelStyle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		PaddingLeft(4)
	builder.WriteString(cmdLabelStyle.Render("Command:") + "\n")

	commandBoxStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		MarginLeft(4).
		Width(width - 12)

	wrappedCmd := wrapText("$ "+command, width-16)
	builder.WriteString(commandBoxStyle.Render(wrappedCmd) + "\n\n")

	return builder.String()
}

// renderOutputBlock renders output in a styled box
func renderOutputBlock(output string, width int, maxLines int) string {
	if output == "" {
		return ""
	}

	var builder strings.Builder

	outputLabelStyle := lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true).
		PaddingLeft(4)

	builder.WriteString(outputLabelStyle.Render("Output:") + "\n")

	outputBoxStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0, 1).
		MarginLeft(4).
		Width(width - 12)

	truncatedOutput := truncateOutput(output, maxLines)
	builder.WriteString(outputBoxStyle.Render(truncatedOutput) + "\n\n")

	return builder.String()
}

// renderErrorBlock renders error in a styled box
func renderErrorBlock(errorMsg string, width int, maxLines int) string {
	if errorMsg == "" {
		return ""
	}

	var builder strings.Builder

	errorLabelStyle := lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true).
		PaddingLeft(4)

	builder.WriteString(errorLabelStyle.Render("Error:") + "\n")

	errorBoxStyle := lipgloss.NewStyle().
		Foreground(colorError).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorError).
		Padding(0, 1).
		MarginLeft(4).
		Width(width - 12)

	truncatedError := truncateOutput(errorMsg, maxLines)
	builder.WriteString(errorBoxStyle.Render(truncatedError) + "\n\n")

	return builder.String()
}

// renderStatusBadge renders a status badge for a step
// Handles both standard status codes ("success") and emoji statuses ("✅ Success")
func renderStatusBadge(status string, isCurrent bool, isExecuteMode bool) string {
	// Normalize status to handle both formats
	normalizedStatus := normalizeStatus(status)
	
	var badge string

	switch normalizedStatus {
	case "success":
		badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(colorSuccess).
			Padding(0, 1).
			Bold(true).
			Render("✓ DONE")
	case "error":
		badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(colorError).
			Padding(0, 1).
			Bold(true).
			Render("✗ ERROR")
	case "skipped":
		badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(colorWarning).
			Padding(0, 1).
			Bold(true).
			Render("⊘ SKIPPED")
	case "timeout":
		badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(colorWarning).
			Padding(0, 1).
			Bold(true).
			Render("⏰ TIMEOUT")
	default:
		if isCurrent && isExecuteMode {
			badge = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(colorAccent).
				Padding(0, 1).
				Bold(true).
				Render("▶ CURRENT")
		} else {
			badge = lipgloss.NewStyle().
				Foreground(colorFaint).
				Background(lipgloss.Color("235")).
				Padding(0, 1).
				Render("○ PENDING")
		}
	}

	// Add viewing indicator for log mode
	if isCurrent && !isExecuteMode {
		return "  " + badge + " " + lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			Render("◀ VIEWING")
	}

	return "  " + badge
}

// renderTitleHeader renders a centered title with divider
func renderTitleHeader(title string, width int) string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(width - 4).
		MarginBottom(1)

	builder.WriteString(titleStyle.Render(title) + "\n")

	// Horizontal divider
	divider := strings.Repeat("─", width-4)
	dividerStyle := lipgloss.NewStyle().
		Foreground(colorBorder)
	builder.WriteString(dividerStyle.Render(divider) + "\n\n")

	return builder.String()
}

// renderProgressBar renders a progress bar with label
func renderProgressBar(completed, total int, width int) string {
	var builder strings.Builder

	progressPercent := 0
	if total > 0 {
		progressPercent = (completed * 100) / total
	}

	// Create visual progress bar
	barWidth := 40
	if width < 80 {
		barWidth = width - 40
		if barWidth < 10 {
			barWidth = 10
		}
	}
	filledWidth := (progressPercent * barWidth) / 100
	progressBar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

	progressStyle := lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true)

	progressLabel := fmt.Sprintf("Progress: %d/%d (%d%%)", completed, total, progressPercent)
	builder.WriteString(progressStyle.Render(progressLabel) + "\n")
	builder.WriteString(progressBar + "\n\n")

	return builder.String()
}

// renderStepSeparator renders a separator line between steps
func renderStepSeparator(width int) string {
	separator := lipgloss.NewStyle().
		Foreground(colorFaint).
		Render(strings.Repeat("─", width-4))
	return separator + "\n\n"
}

// normalizeStatus converts emoji status to standard status code
// Handles both "success" and "✅ Success" formats
func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	
	if strings.Contains(status, "✅") || strings.Contains(status, "Success") {
		return "success"
	}
	if strings.Contains(status, "❌") || strings.Contains(status, "Error") {
		return "error"
	}
	if strings.Contains(status, "⏭️") || strings.Contains(status, "Skipped") {
		return "skipped"
	}
	if strings.Contains(status, "⏰") || strings.Contains(status, "Timeout") {
		return "timeout"
	}
	
	return status // Return as-is if no match
}

// extractSOPName extracts the SOP name from a file path
// Example: "/home/user/.opsy/sops/infra/postgres-backup.md" → "postgres-backup"
func extractSOPName(sopPath string) string {
	if sopPath == "" {
		return "Unknown"
	}
	
	// Get base filename
	base := filepath.Base(sopPath)
	
	// Remove .md extension
	name := strings.TrimSuffix(base, ".md")
	
	return name
}
