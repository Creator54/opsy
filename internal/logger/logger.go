package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"opsy/internal/config"
	"opsy/internal/types"
)

func init() {
	// Ensure time package is considered used
	_ = time.Now()
}

// Logger handles logging of SOP executions
type Logger struct {
	logDirectory string
}

// NewLogger creates a new logger with the default or configured log directory
func NewLogger() (*Logger, error) {
	cfg := config.GetConfig()
	
	// Create the log directory if it doesn't exist
	if err := os.MkdirAll(cfg.LogDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	
	return &Logger{
		logDirectory: cfg.LogDirectory,
	}, nil
}

// LogExecution saves the execution results of an SOP to a log file
func (l *Logger) LogExecution(execution types.SOPExecution) (string, error) {
	// Create the daily subdirectory (e.g., 09-10-2025)
	dateStr := execution.StartedAt.Format("02-01-2006") // DD-MM-YYYY format
	dailyLogDir := filepath.Join(l.logDirectory, dateStr)
	
	// Create subdirectory for the SOP's folder
	sopDir := filepath.Base(filepath.Dir(execution.SOPPath))
	if sopDir == "." || sopDir == "/" {
		sopDir = "default" // Use default if SOP is in root
	}
	logSubDir := filepath.Join(dailyLogDir, sopDir)
	
	if err := os.MkdirAll(logSubDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log subdirectory: %w", err)
	}
	
	// Create filename with timestamp (e.g., deploy-nginx_22-37-14.log.md)
	sopName := strings.TrimSuffix(filepath.Base(execution.SOPPath), ".md")
	timestamp := execution.StartedAt.Format("15-04-05") // HH-MM-SS format
	filename := fmt.Sprintf("%s_%s.log.md", sopName, timestamp)
	logPath := filepath.Join(logSubDir, filename)
	
	// Convert execution to log file format
	logFile := l.executionToLogFile(execution)
	
	// Write the log file
	content := l.formatLogContent(logFile)
	if err := os.WriteFile(logPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write log file: %w", err)
	}
	
	return logPath, nil
}

// executionToLogFile converts an SOP execution to the log file format
func (l *Logger) executionToLogFile(execution types.SOPExecution) types.LogFile {
	logFile := types.LogFile{
		Title:       execution.SOPName,
		SOPRunID:    execution.ID,
		OriginalSOP: execution.SOPPath,
		ExecutedBy:  execution.ExecutedBy,
		StartedAt:   execution.StartedAt,
		EndedAt:     execution.EndedAt,
		Status:      execution.Status,
		Steps:       []types.LogStep{},
	}
	
	// Convert execution steps to log steps
	for _, execStep := range execution.ExecutionLog {
		logStep := types.LogStep{
			StepID:       execStep.StepID,
			Command:      execStep.OriginalStep.Command,
			OriginalStep: execStep.OriginalStep,
		}
		
		if execStep.ExecutionResult != nil {
			logStep.ExecutedAt = execStep.ExecutionResult.ExecutedAt
			logStep.ResultStatus = execStep.ExecutionResult.Status
			logStep.Output = execStep.ExecutionResult.Output
			logStep.ExecutionResult = execStep.ExecutionResult
		}
		
		logFile.Steps = append(logFile.Steps, logStep)
	}
	
	return logFile
}

// formatLogContent formats the log file content according to the PRD specification
func (l *Logger) formatLogContent(logFile types.LogFile) string {
	var content strings.Builder
	
	// Write the header
	content.WriteString(fmt.Sprintf("# %s\n\n", logFile.Title))
	content.WriteString("> **SOP Run ID:** " + logFile.SOPRunID + "  \n")
	content.WriteString("> **Original SOP:** " + logFile.OriginalSOP + "  \n")
	content.WriteString("> **Executed by:** " + logFile.ExecutedBy + "  \n")
	content.WriteString("> **Started at:** " + logFile.StartedAt.Format("2006-01-02 15:04:05") + "  \n")
	content.WriteString("> **Ended at:** " + logFile.EndedAt.Format("2006-01-02 15:04:05") + "  \n")
	
	// Convert status to emoji
	statusEmoji := "✅ Completed Successfully"
	if logFile.Status == "failed" {
		statusEmoji = "❌ Failed"
	} else if logFile.Status == "interrupted" {
		statusEmoji = "⚠️ Interrupted"
	}
	content.WriteString("> **Status:** " + statusEmoji + "\n\n")
	
	// Write each step
	for _, step := range logFile.Steps {
		content.WriteString(fmt.Sprintf("## Step %d: %s\n", step.StepID, step.OriginalStep.Title))
		
		// Write the command
		content.WriteString("```bash\n")
		content.WriteString(step.Command + "\n")
		content.WriteString("```\n\n")
		
		if step.ExecutionResult != nil {
			// Write execution details
			content.WriteString("> **Executed:** " + step.ExecutedAt.Format("2006-01-02 15:04:05") + "  \n")
			
			// Convert result status to emoji
			resultEmoji := "✅ Success"
			if step.ResultStatus == "error" {
				resultEmoji = "❌ Error"
			} else if step.ResultStatus == "timeout" {
				resultEmoji = "⏰ Timeout"
			} else if step.ResultStatus == "skipped" {
				resultEmoji = "⏭️ Skipped"
			}
			content.WriteString("> **Result:** " + resultEmoji + "  \n")
			
			// Write output if available
			if step.Output != "" {
				content.WriteString("> **Output:**\n")
				content.WriteString("> ```\n")
				// Ensure output is properly formatted with > prefix for each line
				for _, line := range strings.Split(step.Output, "\n") {
					content.WriteString("> " + line + "\n")
				}
				content.WriteString("> ```\n")
			}
		}
		
		content.WriteString("\n")
	}
	
	return content.String()
}

// GetLogDirectory returns the logger's log directory
func (l *Logger) GetLogDirectory() string {
	return l.logDirectory
}