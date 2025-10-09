package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"opsy/internal/types"
)

// saveExecutionLog saves the current execution log
func (m model) saveExecutionLog() tea.Cmd {
	return func() tea.Msg {
		execution := types.SOPExecution{
			ID:           "test-run",
			SOPName:      m.sop.Title,
			SOPPath:      m.sop.Path,
			ExecutedBy:   "user",
			Status:       "completed",
			ExecutionLog: []types.ExecutionStep{},
		}

		// Convert m.steps to execution log
		for i, step := range m.steps {
			execution.ExecutionLog = append(execution.ExecutionLog, types.ExecutionStep{
				StepID:       step.ID,
				OriginalStep: m.sop.Steps[i],
				ExecutionResult: &types.ExecutionResult{
					Status: step.Status,
					Output: step.Output,
				},
			})
		}

		logPath, err := m.logger.LogExecution(execution)
		if err != nil {
			return enterModeMsg{
				mode:   modeExecute,
				status: fmt.Sprintf("Error saving log: %v", err),
			}
		}
		return enterModeMsg{
			mode:   modeExecute,
			status: fmt.Sprintf("Log saved to: %s", logPath),
		}
	}
}
