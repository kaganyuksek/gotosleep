package shutdown

import (
	"fmt"
	"os/exec"
	"strconv"
)

// WindowsExecutor implements Executor for Windows
type WindowsExecutor struct{}

// Schedule schedules a shutdown on Windows
func (e *WindowsExecutor) Schedule(minutes int, dryRun bool) (string, error) {
	seconds := minutes * 60
	command := fmt.Sprintf("shutdown.exe /s /t %d", seconds)

	if dryRun {
		return command, nil
	}

	cmd := exec.Command("shutdown.exe", "/s", "/t", strconv.Itoa(seconds))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return command, fmt.Errorf("failed to schedule shutdown: %v, output: %s", err, string(output))
	}

	return command, nil
}

// Cancel cancels a scheduled shutdown on Windows
func (e *WindowsExecutor) Cancel(dryRun bool) error {
	if dryRun {
		return nil
	}

	cmd := exec.Command("shutdown.exe", "/a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to cancel shutdown: %v, output: %s", err, string(output))
	}

	return nil
}

// GetOS returns the OS name
func (e *WindowsExecutor) GetOS() string {
	return "windows"
}
