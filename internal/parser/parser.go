package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"opsy/internal/types"
)

// ParseSOP parses a markdown file and extracts executable command blocks
func ParseSOP(filePath string) (*types.SOP, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	sop := &types.SOP{
		Name: filePath,
		Path: filePath,
		Steps: []types.Step{},
	}

	// Extract title from the first H1 header if available
	title := extractTitle(lines)
	if title != "" {
		sop.Title = title
	}

	stepID := 1
	inCodeBlock := false
	var currentCodeBlock strings.Builder
	currentCodeType := ""
	currentStepLineNumber := 0

	// Pattern to match code fences like ```bash or ```sh
	codeFenceStart := regexp.MustCompile("^```(\\w+)")
	
	for i, line := range lines {
		// Check if we're starting a code block
		startMatches := codeFenceStart.FindStringSubmatch(line)
		if len(startMatches) > 1 {
			lang := strings.ToLower(startMatches[1])
			// Only process bash/shell code blocks
			if lang == "bash" || lang == "sh" || lang == "shell" {
				inCodeBlock = true
				currentCodeType = lang
				currentCodeBlock.Reset()
				currentStepLineNumber = i + 1 // Line numbers start from 1
				continue
			}
		}

		// If in a code block and encounter closing fence
		if inCodeBlock && strings.HasPrefix(line, "```") {
			inCodeBlock = false
			command := strings.TrimSpace(currentCodeBlock.String())
			if command != "" {
				// Find the description/title for this step
				stepDescription := findStepDescription(lines, currentStepLineNumber)
				
				step := types.Step{
					ID:          stepID,
					Title:       extractTitleFromCommand(command), // Use first few words as title
					Description: stepDescription,
					Command:     command,
					CommandType: currentCodeType,
					LineNumber:  currentStepLineNumber,
				}
				sop.Steps = append(sop.Steps, step)
				stepID++
			}
			continue
		}

		// If in a code block, add line to current code
		if inCodeBlock {
			currentCodeBlock.WriteString(line)
			currentCodeBlock.WriteString("\n")
			continue
		}
	}

	// Set the modification time
	if fileInfo, err := os.Stat(filePath); err == nil {
		sop.Modified = fileInfo.ModTime()
	}

	return sop, nil
}

// extractTitle extracts the title from the first H1 header
func extractTitle(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
		// Check for Setext-style H1 (underlined with ===)
		if i := indexOf(lines, line); i < len(lines)-1 {
			nextLine := lines[i+1]
			if strings.HasPrefix(nextLine, "=") && isAllEqualSigns(nextLine) {
				return strings.TrimSpace(line)
			}
		}
	}
	return ""
}

func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

func isAllEqualSigns(s string) bool {
	for _, c := range s {
		if c != '=' {
			return false
		}
	}
	return true
}

// findStepDescription finds the text description before the code block
func findStepDescription(lines []string, codeLineNumber int) string {
	// Look backwards from the code block to find the description
	for i := codeLineNumber - 2; i >= 0; i-- { // -2 to skip the code fence line and start from one before the code
		line := lines[i]
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		// If it's a header, return it as the description
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimLeft(line, "# "))
		}
		// If it's not another code block, return it as description
		if !strings.HasPrefix(line, "```") {
			return strings.TrimSpace(line)
		}
		// If we hit another code block, stop searching
		break
	}
	return ""
}

// extractTitleFromCommand creates a title from the command
func extractTitleFromCommand(command string) string {
	cmd := strings.Fields(command)[0] // First word is usually the command
	if len(command) > 50 {
		return cmd + "..." // Truncate if too long
	}
	return cmd
}