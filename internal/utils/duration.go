package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses various duration formats and returns minutes
// Supported formats:
// - "90" -> 90 minutes
// - "90m" -> 90 minutes
// - "1h30m" -> 90 minutes
// - "00:45" -> 45 minutes
// - "1:20" -> 80 minutes (1 hour 20 minutes)
// - "2h" -> 120 minutes
func ParseDuration(input string) (int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Try parsing as plain number (minutes)
	if val, err := strconv.Atoi(input); err == nil {
		if val <= 0 {
			return 0, fmt.Errorf("duration must be positive")
		}
		return val, nil
	}

	// Try parsing HH:MM format
	if matched, _ := regexp.MatchString(`^\d{1,2}:\d{2}$`, input); matched {
		parts := strings.Split(input, ":")
		hours, _ := strconv.Atoi(parts[0])
		mins, _ := strconv.Atoi(parts[1])
		total := hours*60 + mins
		if total <= 0 {
			return 0, fmt.Errorf("duration must be positive")
		}
		return total, nil
	}

	// Try parsing with time.ParseDuration (supports 1h30m, 90m, 2h, etc.)
	// But we need to handle cases without suffixes
	if !strings.ContainsAny(input, "hms") {
		// If it's just a number string, we already handled it
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	duration, err := time.ParseDuration(input)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	minutes := int(duration.Minutes())
	if minutes <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	return minutes, nil
}

// FormatDuration formats minutes into a readable string
func FormatDuration(minutes int) string {
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}

// FormatCountdown formats a duration into HH:MM:SS
func FormatCountdown(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	if totalSeconds < 0 {
		totalSeconds = 0
	}

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// GenerateID generates a simple ID based on timestamp
func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
