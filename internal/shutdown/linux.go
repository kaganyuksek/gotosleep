package shutdown

import (
	"fmt"
	"os/exec"
	"strconv"
)

// LinuxExecutor implements Executor for Linux
type LinuxExecutor struct{}

// Schedule schedules a shutdown on Linux
func (e *LinuxExecutor) Schedule(minutes int, dryRun bool) (string, error) {
	command := fmt.Sprintf("shutdown -h +%d", minutes)

	if dryRun {
		return command, nil
	}

	cmd := exec.Command("shutdown", "-h", "+"+strconv.Itoa(minutes))
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a permission error
		return command, fmt.Errorf("failed to schedule shutdown (may need sudo): %v, output: %s", err, string(output))
	}

	return command, nil
}

// Cancel cancels a scheduled shutdown on Linux
func (e *LinuxExecutor) Cancel(dryRun bool) error {
	if dryRun {
		return nil
	}

	cmd := exec.Command("shutdown", "-c")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to cancel shutdown (may need sudo): %v, output: %s", err, string(output))
	}

	return nil
}

// GetOS returns the OS name
func (e *LinuxExecutor) GetOS() string {
	return "linux"
}
