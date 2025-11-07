package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"opsy/internal/types"
)

// saveExecutionLog saves the current execution log
// Only saves if there are actually executed steps (not just opened for reading)
func (m model) saveExecutionLog() tea.Cmd {
	return func() tea.Msg {
		// Check if any steps have been executed
		hasExecutedSteps := false
		for _, step := range m.steps {
			if step.Status != "" && step.Status != "pending" {
				hasExecutedSteps = true
				break
			}
		}
		
		// Only save log if steps have been executed
		if !hasExecutedSteps {
			// Don't save log for SOPs that were just opened for reading
			// This is normal behavior when users just browse SOPs without executing
			return enterModeMsg{
				mode:   modeExecute,
				status: "",
			}
		}
		
		// Generate a unique ID based on timestamp
		now := time.Now()
		id := fmt.Sprintf("run-%d", now.Unix())
		
		execution := types.SOPExecution{
			ID:        id,
			SOPName:   m.sop.Title,
			SOPPath:   m.sop.Path,
			ExecutedBy: "user", // TODO: Get actual user
			StartedAt: now,
			EndedAt:   now, // For logs, we set both to same time
			Status:    "completed",
			ExecutionLog: []types.ExecutionStep{},
		}

		// Convert m.steps to execution log
		for i, step := range m.steps {
			execution.ExecutionLog = append(execution.ExecutionLog, types.ExecutionStep{
				StepID:       step.ID,
				OriginalStep: m.sop.Steps[i],
				ExecutionResult: &types.ExecutionResult{
					ExecutedAt: step.ExecutedAt,
					Status:     step.Status,
					Output:     step.Output,
				},
			})
		}

		_, err := m.logger.LogExecution(execution)
		if err != nil {
			return enterModeMsg{
				mode:   modeExecute,
				status: fmt.Sprintf("Error saving log: %v", err),
			}
		}
		
		// Silent success - no status message needed for incremental saves
		return enterModeMsg{
			mode:   modeExecute,
			status: "",
		}
	}
}
