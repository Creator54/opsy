package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"opsy/internal/types"
)

// Executor handles the execution of commands from SOP steps
type Executor struct {
	Timeout time.Duration // Maximum time to wait for command execution
}

// NewExecutor creates a new executor with default timeout
func NewExecutor() *Executor {
	return &Executor{
		Timeout: 30 * time.Second, // Default 30 second timeout
	}
}

// ExecuteStep executes a single SOP step and returns the execution result
func (e *Executor) ExecuteStep(step types.Step) (*types.ExecutionResult, error) {
	if step.Command == "" {
		return nil, fmt.Errorf("step has no command to execute")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	// Create the command
	cmd := exec.CommandContext(ctx, "sh", "-c", step.Command)
	
	// Capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Execute the command
	err := cmd.Run()
	endTime := time.Now()

	// Build the output string from both stdout and stderr
	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderrBuf.String()
	}

	result := &types.ExecutionResult{
		ExecutedAt: endTime,
		Output:     strings.TrimSpace(output),
	}

	// Determine status based on execution result
	if ctx.Err() == context.DeadlineExceeded {
		result.Status = "timeout"
		result.Error = "Command timed out"
		result.ExitCode = -1
	} else if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1 // Generic error code
		}
	} else {
		result.Status = "success"
		result.ExitCode = 0
	}

	return result, nil
}

// ExecuteWithInput executes a command with provided input
func (e *Executor) ExecuteWithInput(step types.Step, input string) (*types.ExecutionResult, error) {
	if step.Command == "" {
		return nil, fmt.Errorf("step has no command to execute")
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", step.Command)
	
	var stdinBuf bytes.Buffer
	var stdoutBuf, stderrBuf bytes.Buffer
	
	// Set up input if provided
	if input != "" {
		stdinBuf.WriteString(input)
		cmd.Stdin = &stdinBuf
	}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	endTime := time.Now()

	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderrBuf.String()
	}

	result := &types.ExecutionResult{
		ExecutedAt: endTime,
		Output:     strings.TrimSpace(output),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Status = "timeout"
		result.Error = "Command timed out"
		result.ExitCode = -1
	} else if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
	} else {
		result.Status = "success"
		result.ExitCode = 0
	}

	return result, nil
}

// ValidateCommand checks if a command is safe to execute
// This is a basic safety check - more sophisticated validation can be added
func (e *Executor) ValidateCommand(command string) error {
	// Check for potentially dangerous commands
	blockedCommands := []string{
		"rm -rf /",      // Delete entire filesystem
		"rm -rf /*",     // Delete entire root
		":(){:|:&};:",   // Fork bomb
		"mkfs.",         // File system creation/formatting commands
	}
	
	for _, blocked := range blockedCommands {
		if strings.Contains(command, blocked) {
			return fmt.Errorf("command contains potentially dangerous pattern: %s", blocked)
		}
	}
	
	return nil
}

// ExecuteCommand executes a command string directly and returns the result
func (e *Executor) ExecuteCommand(command string) (*types.ExecutionResult, error) {
	// Create a temporary step to use with ExecuteStep
	step := types.Step{
		Command: command,
	}
	
	return e.ExecuteStep(step)
}