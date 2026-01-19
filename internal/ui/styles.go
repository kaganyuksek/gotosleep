package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#04B575")
	errorColor     = lipgloss.Color("#FF6B6B")
	warningColor   = lipgloss.Color("#FFD93D")
	textColor      = lipgloss.Color("#FAFAFA")
	dimColor       = lipgloss.Color("#7D7D7D")
	borderColor    = lipgloss.Color("#383838")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)

	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	BigTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true)

	StatusActiveStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 2).
			MarginRight(1)

	ButtonSecondaryStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(dimColor).
				Padding(0, 2).
				MarginRight(1)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Padding(0, 2).
				MarginRight(1)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(textColor).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			Width(30)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(textColor).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(secondaryColor).
				Padding(0, 1).
				Width(30)

	// Countdown styles
	CountdownStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginTop(1).
			MarginBottom(1)

	BigCountdownStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				Align(lipgloss.Center).
				MarginTop(2).
				MarginBottom(2)

	// Progress bar styles
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)

	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(dimColor)

	// Help text styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Padding(0, 2)

	// Preset styles
	PresetStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1).
			MarginRight(1)

	PresetKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)
)

// RenderProgressBar renders a progress bar
func RenderProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int((percent / 100.0) * float64(width))
	empty := width - filled

	var bar strings.Builder
	for i := 0; i < filled; i++ {
		bar.WriteString("█")
	}
	for i := 0; i < empty; i++ {
		bar.WriteString("░")
	}

	return ProgressBarStyle.Render(bar.String())
}
