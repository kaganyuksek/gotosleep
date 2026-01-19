package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaganyuksek/gotosleep/internal/config"
	"github.com/kaganyuksek/gotosleep/internal/utils"
)

// TickMsg is sent every second for countdown updates
type TickMsg time.Time

// ActiveModel represents the active countdown screen
type ActiveModel struct {
	config    *config.Config
	width     int
	height    int
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
}

// NewActiveModel creates a new active model
func NewActiveModel(cfg *config.Config) ActiveModel {
	var startTime, endTime time.Time
	var duration time.Duration

	if cfg.ActiveJob != nil {
		startTime = cfg.ActiveJob.StartTime
		endTime = cfg.ActiveJob.EndTime
		duration = endTime.Sub(startTime)
	}

	return ActiveModel{
		config:    cfg,
		startTime: startTime,
		endTime:   endTime,
		duration:  duration,
	}
}

// Init initializes the active model
func (m ActiveModel) Init() tea.Cmd {
	return tick()
}

// tick returns a command that sends a TickMsg every second
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages for the active screen
func (m ActiveModel) Update(msg tea.Msg) (ActiveModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case TickMsg:
		// Continue ticking
		return m, tick()
	}

	return m, nil
}

// View renders the active countdown screen
func (m ActiveModel) View() string {
	var s strings.Builder

	// Calculate remaining time
	now := time.Now()
	remaining := m.endTime.Sub(now)
	if remaining < 0 {
		remaining = 0
	}

	// Calculate progress percentage
	elapsed := now.Sub(m.startTime)
	if elapsed < 0 {
		elapsed = 0
	}
	progress := (float64(elapsed) / float64(m.duration)) * 100
	if progress > 100 {
		progress = 100
	}

	// Title
	title := BigTitleStyle.Render("Shutting down in")
	s.WriteString(title + "\n\n")

	// Calculate dynamic widths based on content area (inside the border)
	contentAreaWidth := max(m.width-8, 40) // Full width usage
	progressBarWidth := max(contentAreaWidth-8, 30)

	// Big countdown
	countdown := utils.FormatCountdown(remaining)
	bigCountdown := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true).
		Align(lipgloss.Center).
		Width(contentAreaWidth).
		Render(countdown)

	// Make it really big
	s.WriteString(lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Render(lipgloss.NewStyle().
			Padding(1).
			Render(bigCountdown)))
	s.WriteString("\n")

	// Progress bar
	progressBar := RenderProgressBar(progress, progressBarWidth)
	progressText := fmt.Sprintf("  %.0f%%", progress)
	s.WriteString(lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(contentAreaWidth).
		Render(progressBar + progressText))
	s.WriteString("\n\n")

	// Scheduled time info
	info := fmt.Sprintf("Started: %s  →  Scheduled: %s",
		m.startTime.Format("15:04:05"),
		m.endTime.Format("15:04:05"))
	s.WriteString(StatusStyle.Render(info) + "\n\n")

	// Actions
	help := ""
	help += KeyStyle.Render("c") + " Cancel   "
	help += KeyStyle.Render("e") + " Edit   "
	help += KeyStyle.Render("h") + " History   "
	help += KeyStyle.Render("Esc") + " Back"
	s.WriteString(HelpStyle.Render(help))

	// Wrap in box with responsive width
	contentWidth := max(m.width-2, 40)
	content := BaseStyle.Width(contentWidth).Render(s.String())
	return content
}

// Refresh updates the active model with latest config
func (m *ActiveModel) Refresh(cfg *config.Config) {
	m.config = cfg
	if cfg.ActiveJob != nil {
		m.startTime = cfg.ActiveJob.StartTime
		m.endTime = cfg.ActiveJob.EndTime
		m.duration = m.endTime.Sub(m.startTime)
	}
}
