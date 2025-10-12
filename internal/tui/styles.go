package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	colorPrimary   = lipgloss.Color("213") // Purple
	colorSecondary = lipgloss.Color("170") // Light purple
	colorAccent    = lipgloss.Color("45")  // Cyan
	colorText      = lipgloss.Color("15")  // White
	colorFaint     = lipgloss.Color("240") // Gray
	colorBorder    = lipgloss.Color("242") // Dark gray
	colorSuccess   = lipgloss.Color("42")  // Green
	colorError     = lipgloss.Color("203") // Red
	colorWarning   = lipgloss.Color("220") // Yellow
)

// Base styles
var (
	baseStyle = lipgloss.NewStyle().
			Foreground(colorText)

	headerStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Padding(0, 1).
			Width(80)

	// Item styles for list
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(colorText)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(colorText).
				Background(lipgloss.Color("235")).
				BorderLeft(true).
				BorderLeftForeground(colorAccent)
)
