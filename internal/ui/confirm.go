package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaganyuksek/gotosleep/internal/i18n"
	"github.com/kaganyuksek/gotosleep/internal/utils"
)

// ConfirmModel represents a confirmation dialog
type ConfirmModel struct {
	message   string
	minutes   int
	dryRun    bool
	width     int
	height    int
	confirmed bool
	cancelled bool
}

// NewConfirmModel creates a new confirm model with dry-run setting from config
func NewConfirmModel(minutes int, dryRunDefault bool) ConfirmModel {
	return ConfirmModel{
		message:   i18n.T("confirm.title"),
		minutes:   minutes,
		dryRun:    dryRunDefault,
		confirmed: false,
		cancelled: false,
	}
}

// Init initializes the confirm model
func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the confirm dialog
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			return m, nil
		case "n", "N", "esc":
			m.cancelled = true
			return m, nil
		case "d", "D":
			m.dryRun = !m.dryRun
			return m, nil
		}
	}

	return m, nil
}

// View renders the confirm dialog
func (m ConfirmModel) View() string {
	var s strings.Builder

	// Title
	title := TitleStyle.Render(m.message)
	s.WriteString(title + "\n\n")

	// Message
	durationStr := utils.FormatDuration(m.minutes)
	msg := fmt.Sprintf("%s %s?", i18n.T("confirm.message"), durationStr)
	s.WriteString(lipgloss.NewStyle().Bold(true).Render(msg) + "\n\n")

	// Options
	yesBtn := KeyStyle.Render("[Y]") + " Yes   "
	noBtn := KeyStyle.Render("[N]") + " No   "
	s.WriteString(yesBtn + noBtn + "\n\n")

	// Dry-run toggle
	dryRunLabel := i18n.T("confirm.dry_run") + ": "
	if m.dryRun {
		dryRunLabel += lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("✓ ON")
	} else {
		dryRunLabel += lipgloss.NewStyle().Foreground(lipgloss.Color("#7D7D7D")).Render("✗ OFF")
	}
	dryRunLabel += "  " + KeyStyle.Render("[D]") + " " + i18n.T("actions.toggle")
	s.WriteString(dryRunLabel + "\n")

	if m.dryRun {
		s.WriteString(WarningStyle.Render(i18n.T("confirm.dry_run_help")) + "\n")
	}

	// Wrap in box with responsive width
	contentWidth := max(m.width-2, 40)
	content := BaseStyle.Width(contentWidth).Render(s.String())
	return content
}

// IsConfirmed returns true if the user confirmed
func (m ConfirmModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns true if the user cancelled
func (m ConfirmModel) IsCancelled() bool {
	return m.cancelled
}

// IsDryRun returns true if dry-run mode is enabled
func (m ConfirmModel) IsDryRun() bool {
	return m.dryRun
}

// Reset resets the confirm state
func (m *ConfirmModel) Reset() {
	m.confirmed = false
	m.cancelled = false
}
