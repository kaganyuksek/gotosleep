package shutdown

import (
	"runtime"
	"time"
)

// Executor represents a shutdown command executor
type Executor interface {
	Schedule(minutes int, dryRun bool) (string, error)
	Cancel(dryRun bool) error
	GetOS() string
}

// NewExecutor creates a new executor based on the current OS
func NewExecutor() Executor {
	switch runtime.GOOS {
	case "windows":
		return &WindowsExecutor{}
	case "linux":
		return &LinuxExecutor{}
	case "darwin":
		return &DarwinExecutor{}
	default:
		return &WindowsExecutor{} // Default to Windows for unknown OS
	}
}

// JobInfo represents information about a scheduled job
type JobInfo struct {
	StartTime   time.Time
	EndTime     time.Time
	DurationSec int
	Command     string
}

// CalculateJobInfo calculates job timing information
func CalculateJobInfo(minutes int) JobInfo {
	now := time.Now()
	duration := time.Duration(minutes) * time.Minute
	endTime := now.Add(duration)

	return JobInfo{
		StartTime:   now,
		EndTime:     endTime,
		DurationSec: minutes * 60,
	}
}
