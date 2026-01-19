package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaganyuksek/gotosleep/internal/config"
)

// SettingsModel represents the settings screen
type SettingsModel struct {
	config       *config.Config
	width        int
	height       int
	selectedItem int
	editing      bool
	input        textinput.Model
	err          string
}

// NewSettingsModel creates a new settings model
func NewSettingsModel(cfg *config.Config) SettingsModel {
	ti := textinput.New()
	ti.CharLimit = 20
	ti.Width = 20

	return SettingsModel{
		config:       cfg,
		selectedItem: 0,
		editing:      false,
		input:        ti,
	}
}

// Init initializes the settings model
func (m SettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the settings screen
func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				// Save the edit
				m.editing = false
				m.err = ""
				return m, nil
			case "esc":
				// Cancel edit
				m.editing = false
				m.input.SetValue("")
				m.err = ""
				return m, nil
			}
			// Update text input
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		// Navigation
		itemCount := 2 + len(m.config.Presets) // confirm, dry-run, + presets

		switch msg.String() {
		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
				m.err = ""
			}

		case "down", "j":
			if m.selectedItem < itemCount-1 {
				m.selectedItem++
				m.err = ""
			}

		case "enter", " ":
			// Toggle settings
			if m.selectedItem == 0 {
				m.config.Settings.Confirm = !m.config.Settings.Confirm
			} else if m.selectedItem == 1 {
				m.config.Settings.DryRunDefault = !m.config.Settings.DryRunDefault
			}
			m.err = ""
		}
	}

	return m, nil
}

// View renders the settings screen
func (m SettingsModel) View() string {
	var s strings.Builder

	// Title
	title := BigTitleStyle.Render("Settings")
	s.WriteString(title + "\n\n")

	// Settings items
	items := []struct {
		label string
		value string
	}{
		{"Always confirm before shutdown", m.formatBool(m.config.Settings.Confirm)},
		{"Dry-run by default", m.formatBool(m.config.Settings.DryRunDefault)},
	}

	for i, item := range items {
		line := fmt.Sprintf("%s: %s", item.label, item.value)

		if i == m.selectedItem && !m.editing {
			line = ListItemSelectedStyle.Render("▶ " + line)
		} else {
			line = ListItemStyle.Render("  " + line)
		}

		s.WriteString(line + "\n")
	}

	s.WriteString("\n")
	s.WriteString(TitleStyle.Render("Presets") + "\n")

	// Display presets
	for i, preset := range m.config.Presets {
		line := fmt.Sprintf("[%d] %s", i+1, preset.Label)
		s.WriteString(ListItemStyle.Render("  "+line) + "\n")
	}

	s.WriteString("\n")

	// Error message
	if m.err != "" {
		s.WriteString(ErrorStyle.Render("Error: "+m.err) + "\n\n")
	}

	// Actions
	help := ""
	help += KeyStyle.Render("↑↓") + " Navigate   "
	help += KeyStyle.Render("Space/Enter") + " Toggle   "
	help += KeyStyle.Render("Esc") + " Back"
	s.WriteString(HelpStyle.Render(help))

	// Wrap in box with responsive width
	contentWidth := max(m.width-2, 40)
	content := BaseStyle.Width(contentWidth).Render(s.String())
	return content
}

// formatBool formats a boolean value with color
func (m SettingsModel) formatBool(value bool) string {
	if value {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true).Render("ON")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#7D7D7D")).Render("OFF")
}

// Refresh updates the settings model with latest config
func (m *SettingsModel) Refresh(cfg *config.Config) {
	m.config = cfg
}
