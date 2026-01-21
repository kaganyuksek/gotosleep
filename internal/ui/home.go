package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaganyuksek/gotosleep/internal/config"
	"github.com/kaganyuksek/gotosleep/internal/i18n"
	"github.com/kaganyuksek/gotosleep/internal/utils"
)

// HomeModel represents the home screen
type HomeModel struct {
	config         *config.Config
	input          textinput.Model
	err            string
	width          int
	height         int
	selectedPreset int
}

// NewHomeModel creates a new home model
func NewHomeModel(cfg *config.Config) HomeModel {
	ti := textinput.New()
	ti.Placeholder = i18n.T("home.placeholder")
	ti.Blur()
	ti.CharLimit = 20
	ti.Width = 30

	return HomeModel{
		config:         cfg,
		input:          ti,
		selectedPreset: -1,
	}
}

// Init initializes the home model
func (m HomeModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the home screen
func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Toggle focus between presets and input
			if m.input.Focused() {
				m.input.Blur()
			} else {
				m.input.Focus()
				m.selectedPreset = -1 // Clear preset selection when focusing input
			}
			return m, textinput.Blink

		case "1", "2", "3", "4", "5", "6":
			// Quick preset selection - only when input is not focused
			if !m.input.Focused() {
				idx := int(msg.String()[0] - '1')
				if idx >= 0 && idx < len(m.config.Presets) {
					m.selectedPreset = idx
					m.err = ""
					return m, nil
				}
			}
			// If input is focused, let it fall through to text input

		case "enter":
			// Either use selected preset or parse input
			if m.selectedPreset >= 0 && m.selectedPreset < len(m.config.Presets) {
				return m, nil // Signal to parent that we want to start
			}

			if m.input.Value() != "" {
				_, err := utils.ParseDuration(m.input.Value())
				if err != nil {
					m.err = err.Error()
					return m, nil
				}
				return m, nil // Signal to parent that we want to start
			}

			m.err = i18n.T("home.error_no_duration")
			return m, nil

		case "esc":
			m.selectedPreset = -1
			m.err = ""
			return m, nil
		}
	}

	// Update text input (only when focused)
	if m.input.Focused() {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

// View renders the home screen
func (m HomeModel) View() string {
	var s strings.Builder

	// Title
	title := BigTitleStyle.Render(i18n.T("home.title"))
	s.WriteString(title + "\n\n")

	// Status
	status := i18n.T("home.status") + ": "
	if m.config.ActiveJob != nil {
		status += StatusActiveStyle.Render(i18n.T("home.status_active"))
	} else {
		status += StatusStyle.Render(i18n.T("home.status_inactive"))
	}
	s.WriteString(status + "\n\n")

	// Quick presets
	s.WriteString(TitleStyle.Render(i18n.T("home.quick_presets")+":") + "\n")
	presetLine := ""
	currentLineWidth := 0
	maxLineWidth := max(m.width-8, 80) // Use full available width

	for i, preset := range m.config.Presets {
		key := PresetKeyStyle.Render(fmt.Sprintf("[%d]", i+1))
		label := preset.Label

		var item string
		if m.selectedPreset == i {
			item = key + ButtonActiveStyle.Render(label) + " "
		} else {
			item = key + PresetStyle.Render(label) + " "
		}

		// Estimate width (rough calculation, key + label + space)
		itemWidth := len(fmt.Sprintf("[%d]", i+1)) + len(preset.Label) + 2

		// Check if adding this item would exceed line width
		if currentLineWidth > 0 && currentLineWidth+itemWidth > maxLineWidth {
			s.WriteString(presetLine + "\n")
			presetLine = item
			currentLineWidth = itemWidth
		} else {
			presetLine += item
			currentLineWidth += itemWidth
		}
	}
	if presetLine != "" {
		s.WriteString(presetLine + "\n")
	}

	s.WriteString("\n")

	// Duration input
	s.WriteString(TitleStyle.Render(i18n.T("home.duration")+":") + "\n")
	s.WriteString(m.input.View() + "\n\n")

	// Error message
	if m.err != "" {
		s.WriteString(ErrorStyle.Render(i18n.T("home.error")+": "+m.err) + "\n\n")
	}

	// Actions
	help := ""
	help += KeyStyle.Render(i18n.T("keys.enter")) + " " + i18n.T("actions.start") + "   "
	help += KeyStyle.Render(i18n.T("keys.tab")) + " " + i18n.T("actions.toggle_input") + "   "
	help += KeyStyle.Render(i18n.T("keys.history")) + " " + i18n.T("actions.history") + "   "
	help += KeyStyle.Render(i18n.T("keys.settings")) + " " + i18n.T("actions.settings") + "   "
	if m.config.ActiveJob != nil {
		help += KeyStyle.Render(i18n.T("keys.active")) + " " + i18n.T("actions.active") + "   "
	}
	help += KeyStyle.Render(i18n.T("keys.quit")) + " " + i18n.T("actions.quit")
	s.WriteString(HelpStyle.Render(help))

	// Wrap in box with responsive width
	contentWidth := max(m.width-2, 40)
	content := BaseStyle.Width(contentWidth).Render(s.String())
	return content
}

// GetSelectedDuration returns the selected duration in minutes
func (m HomeModel) GetSelectedDuration() (int, error) {
	if m.selectedPreset >= 0 && m.selectedPreset < len(m.config.Presets) {
		return m.config.Presets[m.selectedPreset].Minutes, nil
	}

	if m.input.Value() != "" {
		return utils.ParseDuration(m.input.Value())
	}

	return 0, fmt.Errorf("no duration selected")
}

// Reset resets the selection
func (m *HomeModel) Reset() {
	m.selectedPreset = -1
	m.input.SetValue("")
	m.err = ""
}
