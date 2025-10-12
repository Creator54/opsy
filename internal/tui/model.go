package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"opsy/internal/config"
	"opsy/internal/types"
)

// ExecutorInterface defines the interface for command execution
type ExecutorInterface interface {
	ExecuteStep(step types.Step) (*types.ExecutionResult, error)
	ValidateCommand(command string) error
}

// LoggerInterface defines the interface for logging
type LoggerInterface interface {
	LogExecution(execution types.SOPExecution) (string, error)
}

// SOPStep represents a step in the SOP with execution state
type SOPStep struct {
	ID          int
	Title       string
	Description string
	Command     string
	Status      string // "pending", "executed", "skipped", "error"
	Output      string
	Error       string
}

// model represents the application state
type model struct {
	// Core state
	width, height int
	mode          string
	quitting      bool

	// Browse mode
	fileList    list.Model
	currentPath string

	// Execute mode
	sop           *types.SOP
	steps         []SOPStep
	currentStep   int
	viewport      viewport.Model
	viewportReady bool

	// Edit mode
	textInput textinput.Model
	textarea  textarea.Model

	// Log mode
	logList        list.Model
	logViewPort    viewport.Model
	logViewReady   bool
	logViewContent string
	logViewPath    string

	// Services
	executor ExecutorInterface
	logger   LoggerInterface
	status   string
}

// NewModel creates a new TUI model
func NewModel(executor ExecutorInterface, logger LoggerInterface) model {
	// Initialize paths
	defaultDir := config.DefaultBaseDirectory()
	os.MkdirAll(defaultDir, 0755)

	// Create file list with initial items
	fileList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 10)
	fileList.Title = ""
	fileList.Styles.Title = headerStyle
	fileList.SetShowStatusBar(false)
	fileList.SetFilteringEnabled(false)
	fileList.SetShowHelp(false)       // Disable list's help
	fileList.SetShowPagination(false) // Disable pagination info

	// Create log list
	logList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 10)
	logList.Title = ""
	logList.Styles.Title = headerStyle
	logList.SetShowStatusBar(false)
	logList.SetFilteringEnabled(false)
	logList.SetShowHelp(false)       // Disable list's help
	logList.SetShowPagination(false) // Disable pagination info

	// Create text input for editing
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	// Initialize the model
	m := model{
		mode:           modeBrowse,
		currentPath:    defaultDir,
		fileList:       fileList,
		logList:        logList,
		logViewReady:   false,
		executor:       executor,
		logger:         logger,
		textInput:      ti,
		status:         "Ready",
		viewportReady:  false,
	}

	// Build initial file list
	m.fileList = m.buildFileList(m.currentPath)

	return m
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}
