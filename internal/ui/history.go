package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaganyuksek/gotosleep/internal/config"
	"github.com/kaganyuksek/gotosleep/internal/utils"
)

// HistoryModel represents the history screen
type HistoryModel struct {
	config       *config.Config
	width        int
	height       int
	selectedItem int
	scrollOffset int
}

// NewHistoryModel creates a new history model
func NewHistoryModel(cfg *config.Config) HistoryModel {
	return HistoryModel{
		config:       cfg,
		selectedItem: 0,
		scrollOffset: 0,
	}
}

// Init initializes the history model
func (m HistoryModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the history screen
func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		historyLen := len(m.config.History)
		if historyLen == 0 {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
				// Adjust scroll if needed
				if m.selectedItem < m.scrollOffset {
					m.scrollOffset = m.selectedItem
				}
			}

		case "down", "j":
			if m.selectedItem < historyLen-1 {
				m.selectedItem++
				// Adjust scroll if needed
				visibleItems := m.height / 2 // Rough estimate
				if m.selectedItem >= m.scrollOffset+visibleItems {
					m.scrollOffset = m.selectedItem - visibleItems + 1
				}
			}

		case "home":
			m.selectedItem = 0
			m.scrollOffset = 0

		case "end":
			m.selectedItem = historyLen - 1
			visibleItems := m.height / 2
			m.scrollOffset = historyLen - visibleItems
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		}
	}

	return m, nil
}

// View renders the history screen
func (m HistoryModel) View() string {
	var s strings.Builder

	// Title
	title := BigTitleStyle.Render("History")
	s.WriteString(title + "\n\n")

	// Check if history is empty
	if len(m.config.History) == 0 {
		s.WriteString(StatusStyle.Render("No history yet") + "\n\n")
	} else {
		// Display history items
		visibleItems := 10 // Show 10 items max
		start := m.scrollOffset
		end := start + visibleItems
		if end > len(m.config.History) {
			end = len(m.config.History)
		}

		for i := start; i < end; i++ {
			h := m.config.History[i]

			// Format date and time
			dateStr := h.CreatedAt.Format("2006-01-02 15:04")

			// Format duration
			durationStr := utils.FormatDuration(h.DurationSeconds / 60)

			// Format scheduled time
			scheduledStr := h.ScheduledFor.Format("15:04")

			// Format status with color
			statusStr := ""
			switch h.Status {
			case "ok":
				statusStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("OK")
			case "cancelled":
				statusStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD93D")).Render("Cancelled")
			case "failed":
				statusStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("Failed")
			case "dry-run":
				statusStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D7D7D")).Render("Dry-run")
			default:
				statusStr = h.Status
			}

			// Format line
			line := fmt.Sprintf("%s  %s  → %s  %s",
				dateStr,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render(durationStr),
				scheduledStr,
				statusStr,
			)

			// Apply selection style
			if i == m.selectedItem {
				line = ListItemSelectedStyle.Render("▶ " + line)
			} else {
				line = ListItemStyle.Render("  " + line)
			}

			s.WriteString(line + "\n")
		}

		// Show scroll indicator if needed
		if len(m.config.History) > visibleItems {
			indicator := fmt.Sprintf("\n%s (%d/%d)",
				StatusStyle.Render("Use ↑↓ to scroll"),
				m.selectedItem+1,
				len(m.config.History),
			)
			s.WriteString(indicator + "\n")
		}
	}

	s.WriteString("\n")

	// Actions
	help := ""
	if len(m.config.History) > 0 {
		help += KeyStyle.Render("Enter") + " Restart   "
		help += KeyStyle.Render("d") + " Delete   "
	}
	help += KeyStyle.Render("Esc") + " Back"
	s.WriteString(HelpStyle.Render(help))

	// Wrap in box with responsive width
	contentWidth := max(m.width-2, 50)
	content := BaseStyle.Width(contentWidth).Render(s.String())
	return content
}

// GetSelectedHistory returns the currently selected history item
func (m HistoryModel) GetSelectedHistory() *config.History {
	if m.selectedItem >= 0 && m.selectedItem < len(m.config.History) {
		return &m.config.History[m.selectedItem]
	}
	return nil
}

// Refresh updates the history model with latest config
func (m *HistoryModel) Refresh(cfg *config.Config) {
	m.config = cfg
	// Adjust selection if history changed
	if m.selectedItem >= len(cfg.History) && len(cfg.History) > 0 {
		m.selectedItem = len(cfg.History) - 1
	}
	if len(cfg.History) == 0 {
		m.selectedItem = 0
		m.scrollOffset = 0
	}
}
