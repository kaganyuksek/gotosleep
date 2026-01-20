package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaganyuksek/gotosleep/internal/config"
	"github.com/kaganyuksek/gotosleep/internal/i18n"
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
	ti.Placeholder = i18n.T("settings.preset_label_placeholder")

	mi := textinput.New()
	mi.CharLimit = 5
	mi.Width = 10
	mi.Placeholder = i18n.T("settings.preset_minutes_placeholder")

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
							m.err = i18n.T("settings.error_label_empty")
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
							m.err = i18n.T("settings.error_minutes_invalid")
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
		itemCount := 3 + len(m.config.Presets) // confirm, dry-run, language, + presets

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
			} else if m.selectedItem == 2 {
				// Toggle language between en and tr
				if m.config.Settings.Language == "en" {
					m.config.Settings.Language = "tr"
				} else {
					m.config.Settings.Language = "en"
				}
			} else if m.selectedItem >= 3 {
				// Edit preset
				presetIndex := m.selectedItem - 3
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
	title := BigTitleStyle.Render(i18n.T("settings.title"))
	s.WriteString(title + "\n\n")

	// Settings items
	items := []struct {
		label string
		value string
	}{
		{i18n.T("settings.confirm_label"), m.formatBool(m.config.Settings.Confirm)},
		{i18n.T("settings.dry_run_label"), m.formatBool(m.config.Settings.DryRunDefault)},
		{i18n.T("settings.language"), m.formatLanguage(m.config.Settings.Language)},
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
	s.WriteString(TitleStyle.Render(i18n.T("settings.presets_title")) + "\n")

	// Display presets
	for i, preset := range m.config.Presets {
		itemIndex := 3 + i
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
		s.WriteString(TitleStyle.Render(i18n.T("settings.edit_preset")) + "\n")

		labelLine := i18n.T("settings.preset_label_placeholder") + ": " + m.input.View()
		if m.editingLabel {
			labelLine = ListItemSelectedStyle.Render("▶ " + labelLine)
		} else {
			labelLine = ListItemStyle.Render("  " + labelLine)
		}
		s.WriteString(labelLine + "\n")

		minutesLine := i18n.T("settings.preset_minutes_placeholder") + ": " + m.minutesInput.View()
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
		s.WriteString(ErrorStyle.Render(i18n.T("home.error")+": "+m.err) + "\n\n")
	}

	// Actions
	help := ""
	if m.editing {
		if m.editingLabel {
			help += KeyStyle.Render(i18n.T("keys.enter")) + " Next   "
		} else {
			help += KeyStyle.Render(i18n.T("keys.enter")) + " Save   "
			help += KeyStyle.Render(i18n.T("keys.tab")) + " " + i18n.T("actions.back") + "   "
		}
		help += KeyStyle.Render(i18n.T("keys.esc")) + " " + i18n.T("actions.cancel")
	} else {
		help += KeyStyle.Render(i18n.T("keys.up")+i18n.T("keys.down")) + " Navigate   "
		help += KeyStyle.Render("Space/"+i18n.T("keys.enter")) + " Toggle/" + i18n.T("actions.edit") + "   "
		help += KeyStyle.Render(i18n.T("keys.esc")) + " " + i18n.T("actions.back")
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

// formatLanguage formats language code with display name
func (m SettingsModel) formatLanguage(lang string) string {
	display := lang
	switch lang {
	case "en":
		display = "English"
	case "tr":
		display = "Türkçe"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render(display)
}

// Refresh updates the settings model with latest config
func (m *SettingsModel) Refresh(cfg *config.Config) {
	m.config = cfg
}
