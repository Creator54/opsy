package tui

import (
	"testing"

	"opsy/internal/types"

	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing
type MockExecutor struct{}

func (m *MockExecutor) ExecuteStep(step types.Step) (*types.ExecutionResult, error) {
	return &types.ExecutionResult{
		Status: "success",
		Output: "Command output",
	}, nil
}

func (m *MockExecutor) ValidateCommand(command string) error {
	return nil
}

type MockLogger struct{}

func (m *MockLogger) LogExecution(execution types.SOPExecution) (string, error) {
	return "/tmp/test.log", nil
}

func TestNewModel(t *testing.T) {
	executor := &MockExecutor{}
	logger := &MockLogger{}
	
	model := NewModel(executor, logger)
	
	assert.Equal(t, modeBrowse, model.mode)
	assert.Equal(t, "Ready", model.status)
	assert.NotNil(t, model.fileList)
}

func TestModelInitialization(t *testing.T) {
	executor := &MockExecutor{}
	logger := &MockLogger{}
	
	model := NewModel(executor, logger)
	
	assert.Equal(t, modeBrowse, model.mode)
	assert.NotNil(t, model.executor)
	assert.NotNil(t, model.logger)
}

func TestGetModeContext(t *testing.T) {
	executor := &MockExecutor{}
	logger := &MockLogger{}
	model := NewModel(executor, logger)
	
	// Test different modes
	model.mode = modeBrowse
	assert.Equal(t, "Browser", model.getModeContext())
	
	model.mode = modeExecute
	assert.Equal(t, "Execution", model.getModeContext())
	
	model.mode = modeLogs
	assert.Equal(t, "Logs", model.getModeContext())
	
	model.mode = modeEdit
	assert.Equal(t, "Edit", model.getModeContext())
}