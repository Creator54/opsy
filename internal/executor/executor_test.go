package executor

import (
	"testing"
	"time"

	"opsy/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestExecuteStep(t *testing.T) {
	executor := NewExecutor()
	
	step := types.Step{
		ID:      1,
		Command: "echo 'hello world'",
	}
	
	result, err := executor.ExecuteStep(step)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, 0, result.ExitCode)
	assert.Equal(t, "hello world", result.Output)
}

func TestExecuteStepWithTimeout(t *testing.T) {
	executor := &Executor{Timeout: 100 * time.Millisecond} // Very short timeout
	
	step := types.Step{
		ID:      1,
		Command: "sleep 1", // This will take longer than our timeout
	}
	
	result, err := executor.ExecuteStep(step)
	
	assert.NoError(t, err) // No error from ExecuteStep, timeout is handled internally
	assert.NotNil(t, result)
	assert.Equal(t, "timeout", result.Status)
	assert.Equal(t, -1, result.ExitCode)
	assert.Contains(t, result.Error, "Command timed out")
}

func TestExecuteStepWithError(t *testing.T) {
	executor := NewExecutor()
	
	step := types.Step{
		ID:      1,
		Command: "exit 1", // Command that exits with error
	}
	
	result, err := executor.ExecuteStep(step)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "error", result.Status)
	assert.Equal(t, 1, result.ExitCode)
}

func TestValidateCommand(t *testing.T) {
	executor := NewExecutor()
	
	// Valid command should pass
	err := executor.ValidateCommand("echo hello")
	assert.NoError(t, err)
	
	// Dangerous command should fail
	err = executor.ValidateCommand("rm -rf /")
	assert.Error(t, err)
	
	err = executor.ValidateCommand(":(){:|:&};:")
	assert.Error(t, err)
}