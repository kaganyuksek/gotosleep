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
	editingLabel bool // true = editing label, false = editing minutes
	input        textinput.Model
	minutesInput textinput.Model
	err          string
}

// NewSettingsModel creates a new settings model
func NewSettingsModel(cfg *config.Config) SettingsModel {
	ti := textinput.New()
	ti.CharLimit = 20
	ti.Width = 20
	ti.Placeholder = "Label (e.g., 15m)"

	mi := textinput.New()
	mi.CharLimit = 5
	mi.Width = 10
	mi.Placeholder = "Minutes"

	return SettingsModel{
		config:       cfg,
		selectedItem: 0,
		editing:      false,
		editingLabel: true,
		input:        ti,
		minutesInput: mi,
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
				// Save the edit based on what we're editing
				presetIndex := m.selectedItem - 2
				if presetIndex >= 0 && presetIndex < len(m.config.Presets) {
					if m.editingLabel {
						// Save label and move to minutes editing
						label := m.input.Value()
						if label == "" {
							m.err = "Label cannot be empty"
							return m, nil
						}
						m.config.Presets[presetIndex].Label = label
						m.editingLabel = false
						m.minutesInput.Focus()
						m.minutesInput.SetValue(fmt.Sprintf("%d", m.config.Presets[presetIndex].Minutes))
						return m, textinput.Blink
					} else {
						// Save minutes and finish editing
						minutes := 0
						_, err := fmt.Sscanf(m.minutesInput.Value(), "%d", &minutes)
						if err != nil || minutes <= 0 {
							m.err = "Invalid minutes value"
							return m, nil
						}
						m.config.Presets[presetIndex].Minutes = minutes
						m.editing = false
						m.editingLabel = true
						m.input.Blur()
						m.minutesInput.Blur()
						m.input.SetValue("")
						m.minutesInput.SetValue("")
						m.err = ""
						return m, nil
					}
				}
				m.editing = false
				m.err = ""
				return m, nil
			case "esc":
				// Cancel edit
				m.editing = false
				m.editingLabel = true
				m.input.Blur()
				m.minutesInput.Blur()
				m.input.SetValue("")
				m.minutesInput.SetValue("")
				m.err = ""
				return m, nil
			case "tab":
				// Switch between label and minutes input
				if !m.editingLabel {
					m.editingLabel = true
					m.input.Focus()
					m.minutesInput.Blur()
					return m, textinput.Blink
				}
			}
			// Update appropriate text input
			if m.editingLabel {
				m.input, cmd = m.input.Update(msg)
			} else {
				m.minutesInput, cmd = m.minutesInput.Update(msg)
			}
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
			// Toggle settings or edit presets
			if m.selectedItem == 0 {
				m.config.Settings.Confirm = !m.config.Settings.Confirm
			} else if m.selectedItem == 1 {
				m.config.Settings.DryRunDefault = !m.config.Settings.DryRunDefault
			} else if m.selectedItem >= 2 {
				// Edit preset
				presetIndex := m.selectedItem - 2
				if presetIndex < len(m.config.Presets) {
					m.editing = true
					m.editingLabel = true
					m.input.Focus()
					m.input.SetValue(m.config.Presets[presetIndex].Label)
					return m, textinput.Blink
				}
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
		itemIndex := 2 + i
		line := fmt.Sprintf("%s → %d min", preset.Label, preset.Minutes)

		if itemIndex == m.selectedItem && !m.editing {
			line = ListItemSelectedStyle.Render("▶ " + line)
		} else {
			line = ListItemStyle.Render("  " + line)
		}

		s.WriteString(line + "\n")
	}

	// Show edit form if editing
	if m.editing {
		s.WriteString("\n")
		s.WriteString(TitleStyle.Render("Edit Preset") + "\n")

		labelLine := "Label: " + m.input.View()
		if m.editingLabel {
			labelLine = ListItemSelectedStyle.Render("▶ " + labelLine)
		} else {
			labelLine = ListItemStyle.Render("  " + labelLine)
		}
		s.WriteString(labelLine + "\n")

		minutesLine := "Minutes: " + m.minutesInput.View()
		if !m.editingLabel {
			minutesLine = ListItemSelectedStyle.Render("▶ " + minutesLine)
		} else {
			minutesLine = ListItemStyle.Render("  " + minutesLine)
		}
		s.WriteString(minutesLine + "\n")
	}

	s.WriteString("\n")

	// Error message
	if m.err != "" {
		s.WriteString(ErrorStyle.Render("Error: "+m.err) + "\n\n")
	}

	// Actions
	help := ""
	if m.editing {
		if m.editingLabel {
			help += KeyStyle.Render("Enter") + " Next   "
		} else {
			help += KeyStyle.Render("Enter") + " Save   "
			help += KeyStyle.Render("Tab") + " Back   "
		}
		help += KeyStyle.Render("Esc") + " Cancel"
	} else {
		help += KeyStyle.Render("↑↓") + " Navigate   "
		help += KeyStyle.Render("Space/Enter") + " Toggle/Edit   "
		help += KeyStyle.Render("Esc") + " Back"
	}
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
