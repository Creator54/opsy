package tui

// Mode constants
const (
	modeBrowse  = "browse"
	modeExecute = "execute"
	modeLogs    = "logs"
	modeEdit    = "edit"
)

// UI layout constants
const (
	headerHeight  = 1
	newlineHeight = 1
	spacingHeight = 2
	helpBarHeight = 1
	
	// Total UI chrome height
	uiChromeHeight = headerHeight + newlineHeight // For browse/logs modes
	executeUIChromeHeight = headerHeight + newlineHeight + spacingHeight + helpBarHeight // For execute mode
	
	// Minimum content height
	minContentHeight = 10
)

// Status constants
const (
	statusPending  = "pending"
	statusSuccess  = "success"
	statusError    = "error"
	statusSkipped  = "skipped"
	statusExecuted = "executed"
)
