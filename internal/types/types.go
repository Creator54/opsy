package types

import (
	"time"
)

// SOP represents a Standard Operating Procedure
type SOP struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Steps       []Step     `json:"steps"`
	Modified    time.Time  `json:"modified"`
}

// Step represents a single step in an SOP
type Step struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command"`      // The actual command to execute
	CommandType string `json:"command_type"` // "bash", "shell", etc.
	Executed    bool   `json:"executed"`
	Result      *ExecutionResult `json:"result,omitempty"`
	LineNumber  int    `json:"line_number"`  // Line number in the original markdown file
}

// ExecutionResult holds the result of executing a command
type ExecutionResult struct {
	ExecutedAt time.Time `json:"executed_at"`
	Status     string    `json:"status"`     // "success", "error", "skipped"
	Output     string    `json:"output"`     // Captured stdout/stderr
	ExitCode   int       `json:"exit_code"`
	Error      string    `json:"error,omitempty"`
}

// SOPExecution represents a single execution run of an SOP
type SOPExecution struct {
	ID            string           `json:"id"`
	SOPName       string           `json:"sop_name"`
	SOPPath       string           `json:"sop_path"`
	ExecutedBy    string           `json:"executed_by"`
	StartedAt     time.Time        `json:"started_at"`
	EndedAt       time.Time        `json:"ended_at"`
	Status        string           `json:"status"` // "completed", "failed", "interrupted"
	ExecutionLog  []ExecutionStep  `json:"execution_log"`
}

// ExecutionStep represents a step in the execution log
type ExecutionStep struct {
	StepID       int              `json:"step_id"`
	OriginalStep Step             `json:"original_step"`
	ExecutionResult *ExecutionResult `json:"execution_result"`
}

// LogFile represents the structure of a log file
type LogFile struct {
	Title       string    `json:"title"`
	SOPRunID    string    `json:"sop_run_id"`
	OriginalSOP string    `json:"original_sop"`
	ExecutedBy  string    `json:"executed_by"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
	Status      string    `json:"status"` // completed status
	Steps       []LogStep `json:"steps"`
}

// LogStep represents a step in the log file
type LogStep struct {
	StepID       int              `json:"step_id"`
	Command      string           `json:"command"`
	ExecutedAt   time.Time        `json:"executed_at"`
	ResultStatus string           `json:"result_status"` // success status
	Output       string           `json:"output"`
	OriginalStep Step             `json:"original_step"`
	ExecutionResult *ExecutionResult `json:"execution_result"`
}