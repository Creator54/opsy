package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"opsy/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestLogExecution(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "opsy-logger-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after test
	
	// Create a logger with the test directory
	logger := &Logger{
		logDirectory: tmpDir,
	}
	
	// Create a sample execution
	execution := types.SOPExecution{
		ID:         "2025-10-09_22-37-14",
		SOPName:    "Test SOP",
		SOPPath:    "/home/user/.opsy/sops/test/test-sop.md",
		ExecutedBy: "testuser",
		StartedAt:  time.Now().Add(-10 * time.Minute),
		EndedAt:    time.Now(),
		Status:     "completed",
		ExecutionLog: []types.ExecutionStep{
			{
				StepID: 1,
				OriginalStep: types.Step{
					ID:    1,
					Title: "Test Step",
					Command: "echo 'hello world'",
				},
				ExecutionResult: &types.ExecutionResult{
					ExecutedAt: time.Now(),
					Status:     "success",
					Output:     "hello world",
					ExitCode:   0,
				},
			},
		},
	}
	
	// Log the execution
	logPath, err := logger.LogExecution(execution)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, logPath)
	
	// Check if the log file was created
	_, err = os.Stat(logPath)
	assert.NoError(t, err)
	
	// Check if the content is correct by reading the file
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	
	// Verify the content contains expected elements
	contentStr := string(content)
	assert.Contains(t, contentStr, "# Test SOP")
	assert.Contains(t, contentStr, "**SOP Run ID:** 2025-10-09_22-37-14")
	assert.Contains(t, contentStr, "**Executed by:** testuser")
	assert.Contains(t, contentStr, "```bash\necho 'hello world'\n```")
	assert.Contains(t, contentStr, "**Result:** ✅ Success")
	assert.Contains(t, contentStr, "hello world")
}

func TestFormatLogContent(t *testing.T) {
	logger := &Logger{}
	
	logFile := types.LogFile{
		Title:       "Test SOP",
		SOPRunID:    "2025-10-09_22-37-14",
		OriginalSOP: "/home/user/.opsy/sops/test/test-sop.md",
		ExecutedBy:  "testuser",
		StartedAt:   time.Date(2025, 10, 9, 22, 37, 14, 0, time.UTC),
		EndedAt:     time.Date(2025, 10, 9, 22, 39, 01, 0, time.UTC),
		Status:      "completed",
		Steps: []types.LogStep{
			{
				StepID:  1,
				Command: "echo 'hello world'",
				OriginalStep: types.Step{
					ID:    1,
					Title: "Test Step",
				},
				ExecutedAt:   time.Date(2025, 10, 9, 22, 37, 21, 0, time.UTC),
				ResultStatus: "success",
				Output:       "hello world",
				ExecutionResult: &types.ExecutionResult{
					ExecutedAt: time.Date(2025, 10, 9, 22, 37, 21, 0, time.UTC),
					Status:     "success",
					Output:     "hello world",
					ExitCode:   0,
				},
			},
		},
	}
	
	content := logger.formatLogContent(logFile)
	
	assert.Contains(t, content, "# Test SOP")
	assert.Contains(t, content, "**SOP Run ID:** 2025-10-09_22-37-14")
	assert.Contains(t, content, "**Executed by:** testuser")
	assert.Contains(t, content, "Started at:** 2025-10-09 22:37:14")
	assert.Contains(t, content, "Ended at:** 2025-10-09 22:39:01")
	assert.Contains(t, content, "## Step 1: Test Step")
	assert.Contains(t, content, "```bash\necho 'hello world'\n```")
	assert.Contains(t, content, "**Result:** ✅ Success")
	assert.Contains(t, content, "hello world")
}

func TestNewLogger(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "opsy-config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after test
	
	// Temporarily override the default log directory
	origLogDir := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origLogDir)
	
	logger, err := NewLogger()
	
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	
	// Check if the default log directory was created
	expectedLogDir := filepath.Join(tmpDir, ".opsy", "logs")
	_, err = os.Stat(expectedLogDir)
	assert.NoError(t, err)
}