package tui

import (
	"fmt"
	"strings"
)

// LogStep represents a parsed step from a log file
type LogStep struct {
	StepNumber  int
	Title       string
	Command     string
	Status      string
	ExecutedAt  string
	Output      string
	HasOutput   bool
}

// LogMetadata represents the header metadata from a log file
type LogMetadata struct {
	Title      string
	RunID      string
	SOPPath    string
	ExecutedBy string
	StartedAt  string
	EndedAt    string
	Status     string
}

// ParseLogFile parses a log markdown file into structured data
func ParseLogFile(content string) (LogMetadata, []LogStep) {
	lines := strings.Split(content, "\n")
	
	metadata := parseLogMetadata(lines)
	steps := parseLogSteps(lines)
	
	return metadata, steps
}

// parseLogMetadata extracts metadata from the log file header
func parseLogMetadata(lines []string) LogMetadata {
	metadata := LogMetadata{}
	
	for i, line := range lines {
		// Stop at first step header
		if strings.HasPrefix(line, "## Step ") {
			break
		}
		
		// Parse title (# Title)
		if strings.HasPrefix(line, "# ") {
			metadata.Title = strings.TrimPrefix(line, "# ")
			continue
		}
		
		// Parse metadata lines (> **Key:** Value)
		if strings.HasPrefix(line, "> **") {
			metaLine := strings.TrimPrefix(line, "> ")
			
			if strings.Contains(metaLine, "**SOP Run ID:**") {
				// Extract value after "**SOP Run ID:**"
				value := strings.TrimPrefix(metaLine, "**SOP Run ID:**")
				metadata.RunID = strings.TrimSpace(value)
			} else if strings.Contains(metaLine, "**Original SOP:**") {
				value := strings.TrimPrefix(metaLine, "**Original SOP:**")
				metadata.SOPPath = strings.TrimSpace(value)
			} else if strings.Contains(metaLine, "**Executed by:**") {
				value := strings.TrimPrefix(metaLine, "**Executed by:**")
				metadata.ExecutedBy = strings.TrimSpace(value)
			} else if strings.Contains(metaLine, "**Started at:**") {
				value := strings.TrimPrefix(metaLine, "**Started at:**")
				metadata.StartedAt = strings.TrimSpace(value)
			} else if strings.Contains(metaLine, "**Ended at:**") {
				value := strings.TrimPrefix(metaLine, "**Ended at:**")
				metadata.EndedAt = strings.TrimSpace(value)
			} else if strings.Contains(metaLine, "**Status:**") {
				value := strings.TrimPrefix(metaLine, "**Status:**")
				metadata.Status = strings.TrimSpace(value)
			}
		}
		
		// Safety check to avoid infinite loop
		if i > 50 {
			break
		}
	}
	
	return metadata
}

// parseLogSteps extracts all steps from the log file
func parseLogSteps(lines []string) []LogStep {
	var steps []LogStep
	var currentStep *LogStep
	inCodeBlock := false
	inOutputBlock := false
	outputBlockStarted := false // Track if we've seen the opening ```
	var outputLines []string
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// Detect step header (## Step N: Title)
		if strings.HasPrefix(line, "## Step ") {
			// Save previous step if exists (including any remaining output)
			if currentStep != nil {
				// Save any remaining output that wasn't closed properly
				if len(outputLines) > 0 && !currentStep.HasOutput {
					currentStep.Output = strings.Join(outputLines, "\n")
					currentStep.HasOutput = true
				}
				steps = append(steps, *currentStep)
			}
			
			// Start new step
			stepHeader := strings.TrimPrefix(line, "## ")
			currentStep = &LogStep{}
			
			// Parse "Step N: Title"
			parts := strings.SplitN(stepHeader, ": ", 2)
			if len(parts) >= 1 {
				// Extract step number from "Step N"
				stepNumStr := strings.TrimPrefix(parts[0], "Step ")
				var stepNum int
				_, _ = fmt.Sscanf(stepNumStr, "%d", &stepNum)
				currentStep.StepNumber = stepNum
			}
			if len(parts) >= 2 {
				currentStep.Title = parts[1]
			}
			
			outputLines = []string{}
			inOutputBlock = false
			outputBlockStarted = false
			continue
		}
		
		if currentStep == nil {
			continue
		}
		
		// Handle code blocks (```bash or ```)
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				inCodeBlock = true
			} else {
				inCodeBlock = false
			}
			continue
		}
		
		// Collect command inside code block
		if inCodeBlock {
			currentStep.Command = line
			continue
		}
		
		// Handle metadata lines (> **Key:** Value)
		if strings.HasPrefix(line, "> **") {
			metaLine := strings.TrimPrefix(line, "> ")
			
			// Check for Output marker
			if strings.Contains(metaLine, "**Output:**") {
				inOutputBlock = true
				outputBlockStarted = false // Haven't seen opening ``` yet
				outputLines = []string{}
				continue
			}
			
			// Parse Executed timestamp
			if strings.Contains(metaLine, "**Executed:**") {
				value := strings.TrimPrefix(metaLine, "**Executed:**")
				currentStep.ExecutedAt = strings.TrimSpace(value)
				continue
			}
			
			// Parse Result status
			if strings.Contains(metaLine, "**Result:**") {
				value := strings.TrimPrefix(metaLine, "**Result:**")
				currentStep.Status = strings.TrimSpace(value)
				continue
			}
		}
		
		// Handle output content when in output block
		if inOutputBlock {
			// Check for ``` markers
			if strings.HasPrefix(line, "> ```") {
				if !outputBlockStarted {
					// This is the opening ``` - start collecting
					outputBlockStarted = true
				} else {
					// This is the closing ``` - save and exit
					if len(outputLines) > 0 {
						currentStep.Output = strings.Join(outputLines, "\n")
						currentStep.HasOutput = true
						outputLines = []string{}
					}
					inOutputBlock = false
					outputBlockStarted = false
				}
				continue
			}
			
			// Only collect lines after we've seen the opening ```
			if outputBlockStarted {
				// Collect output line
				if strings.HasPrefix(line, "> ") {
					// Line with "> " prefix - strip it
					outputLines = append(outputLines, strings.TrimPrefix(line, "> "))
				} else if strings.TrimSpace(line) != "" {
					// Continuation line without "> " prefix (like progress bars)
					outputLines = append(outputLines, line)
				}
			}
			continue
		}
	}
	
	// Save last step
	if currentStep != nil {
		if len(outputLines) > 0 {
			currentStep.Output = strings.Join(outputLines, "\n")
			currentStep.HasOutput = true
		}
		steps = append(steps, *currentStep)
	}
	
	return steps
}
