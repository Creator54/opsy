package tui

import "opsy/internal/types"

// Custom messages for state changes
type browseToDirMsg struct {
	path string
}

type enterModeMsg struct {
	mode   string
	status string
	sop    *types.SOP
	steps  []SOPStep
}

// item represents a file/directory in the browser
type item struct {
	title, desc string
	filePath    string
	isDir       bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
