package shutdown

import (
	"fmt"
	"os/exec"
	"strconv"
)

// DarwinExecutor implements Executor for macOS
type DarwinExecutor struct{}

// Schedule schedules a shutdown on macOS
func (e *DarwinExecutor) Schedule(minutes int, dryRun bool) (string, error) {
	command := fmt.Sprintf("sudo shutdown -h +%d", minutes)

	if dryRun {
		return command, nil
	}

	cmd := exec.Command("sudo", "shutdown", "-h", "+"+strconv.Itoa(minutes))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return command, fmt.Errorf("failed to schedule shutdown (need sudo): %v, output: %s", err, string(output))
	}

	return command, nil
}

// Cancel cancels a scheduled shutdown on macOS
func (e *DarwinExecutor) Cancel(dryRun bool) error {
	if dryRun {
		return nil
	}

	// Try shutdown -c first
	cmd := exec.Command("sudo", "shutdown", "-c")
	err := cmd.Run()
	if err == nil {
		return nil
	}

	// If that fails, try killing the shutdown process
	cmd = exec.Command("sudo", "killall", "shutdown")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to cancel shutdown (need sudo): %v, output: %s", err, string(output))
	}

	return nil
}

// GetOS returns the OS name
func (e *DarwinExecutor) GetOS() string {
	return "darwin"
}
