package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaganyuksek/gotosleep/internal/config"
	"github.com/kaganyuksek/gotosleep/internal/i18n"
	"github.com/kaganyuksek/gotosleep/internal/shutdown"
	"github.com/kaganyuksek/gotosleep/internal/ui"
	"github.com/kaganyuksek/gotosleep/internal/utils"
)

// Screen represents different screens in the app
type Screen int

const (
	ScreenHome Screen = iota
	ScreenConfirm
	ScreenActive
	ScreenHistory
	ScreenSettings
)

// App represents the main application model
type App struct {
	config   *config.Config
	executor shutdown.Executor
	screen   Screen
	home     ui.HomeModel
	confirm  ui.ConfirmModel
	active   ui.ActiveModel
	history  ui.HistoryModel
	settings ui.SettingsModel
	err      string
	quitting bool
	width    int
	height   int
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize i18n with configured language
	if err := i18n.Init(cfg.Settings.Language); err != nil {
		return nil, fmt.Errorf("failed to initialize i18n: %w", err)
	}

	executor := shutdown.NewExecutor()

	return &App{
		config:   cfg,
		executor: executor,
		screen:   ScreenHome,
		home:     ui.NewHomeModel(cfg),
		active:   ui.NewActiveModel(cfg),
		history:  ui.NewHistoryModel(cfg),
		settings: ui.NewSettingsModel(cfg),
	}, nil
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	// Set terminal title
	setTitle := func() tea.Msg {
		fmt.Print("\033]0;Shutdown Timer\007")
		return nil
	}

	// If there's an active job, go to active screen
	if a.config.ActiveJob != nil {
		a.screen = ScreenActive
		return tea.Batch(setTitle, a.active.Init())
	}
	return tea.Batch(setTitle, a.home.Init())
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle global messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Store and propagate size to all screens
		a.width = msg.Width
		a.height = msg.Height
		a.home, _ = a.home.Update(msg)
		a.confirm, _ = a.confirm.Update(msg)
		a.active, _ = a.active.Update(msg)
		a.history, _ = a.history.Update(msg)
		a.settings, _ = a.settings.Update(msg)
		// Force re-render
		return a, tea.ClearScreen

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Check if there's an active job
			if a.config.ActiveJob != nil && a.screen != ScreenActive {
				a.err = "Warning: Active shutdown will not be cancelled"
			}
			a.quitting = true
			return a, tea.Quit
		}
	}

	// Route to appropriate screen
	switch a.screen {
	case ScreenHome:
		return a.updateHome(msg)
	case ScreenConfirm:
		return a.updateConfirm(msg)
	case ScreenActive:
		return a.updateActive(msg)
	case ScreenHistory:
		return a.updateHistory(msg)
	case ScreenSettings:
		return a.updateSettings(msg)
	}

	return a, cmd
}

// View renders the current screen
func (a *App) View() string {
	if a.quitting {
		return ""
	}

	switch a.screen {
	case ScreenHome:
		return a.home.View()
	case ScreenConfirm:
		return a.confirm.View()
	case ScreenActive:
		return a.active.View()
	case ScreenHistory:
		return a.history.View()
	case ScreenSettings:
		return a.settings.View()
	}

	return ""
}

// updateHome handles updates for the home screen
func (a *App) updateHome(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			// Go to history
			a.screen = ScreenHistory
			a.history.Refresh(a.config)
			return a, nil
		case "s":
			// Go to settings
			a.screen = ScreenSettings
			a.settings.Refresh(a.config)
			return a, nil
		case "a":
			// Go to active screen if there's an active job
			if a.config.ActiveJob != nil {
				a.screen = ScreenActive
				a.active.Refresh(a.config)
				return a, a.active.Init()
			}
		case "enter":
			// Try to get duration
			minutes, err := a.home.GetSelectedDuration()
			if err != nil {
				a.home, cmd = a.home.Update(msg)
				return a, cmd
			}

			// Check if confirmation is enabled in settings
			if a.config.Settings.Confirm {
				// Show confirm dialog with DryRunDefault from settings
				a.confirm = ui.NewConfirmModel(minutes, a.config.Settings.DryRunDefault)
				a.screen = ScreenConfirm
				return a, nil
			} else {
				// Skip confirmation and start immediately with DryRunDefault setting
				err := a.startShutdown(minutes, a.config.Settings.DryRunDefault)
				if err != nil {
					a.err = err.Error()
					a.home.Reset()
					return a, nil
				}
				// Go to active screen
				a.screen = ScreenActive
				a.active.Refresh(a.config)
				a.home.Reset()
				return a, a.active.Init()
			}
		}
	}

	a.home, cmd = a.home.Update(msg)
	return a, cmd
}

// updateConfirm handles updates for the confirm dialog
func (a *App) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	a.confirm, cmd = a.confirm.Update(msg)

	// Check if user confirmed or cancelled
	if a.confirm.IsConfirmed() {
		// Start the shutdown
		minutes, _ := a.home.GetSelectedDuration()
		err := a.startShutdown(minutes, a.confirm.IsDryRun())
		if err != nil {
			a.err = err.Error()
			a.screen = ScreenHome
			a.home.Reset()
			return a, nil
		}

		// Go to active screen
		a.screen = ScreenActive
		a.active.Refresh(a.config)
		a.home.Reset()
		return a, a.active.Init()
	}

	if a.confirm.IsCancelled() {
		// Go back to home
		a.screen = ScreenHome
		a.confirm.Reset()
		return a, nil
	}

	return a, cmd
}

// updateActive handles updates for the active screen
func (a *App) updateActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			// Cancel shutdown
			err := a.cancelShutdown()
			if err != nil {
				a.err = err.Error()
			}
			a.screen = ScreenHome
			a.home.Reset()
			return a, nil
		case "e":
			// Edit (cancel and go back to home for new input)
			err := a.cancelShutdown()
			if err != nil {
				a.err = err.Error()
			}
			a.screen = ScreenHome
			a.home.Reset()
			return a, nil
		case "h":
			// Go to history (keep countdown running)
			a.screen = ScreenHistory
			a.history.Refresh(a.config)
			return a, nil
		case "esc":
			// Go back to home
			a.screen = ScreenHome
			return a, nil
		}
	}

	a.active, cmd = a.active.Update(msg)
	return a, cmd
}

// updateHistory handles updates for the history screen
func (a *App) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to previous screen
			if a.config.ActiveJob != nil {
				a.screen = ScreenActive
			} else {
				a.screen = ScreenHome
			}
			return a, nil
		case "enter":
			// Restart selected history item
			selected := a.history.GetSelectedHistory()
			if selected != nil {
				minutes := selected.DurationSeconds / 60

				// Check if confirmation is enabled in settings
				if a.config.Settings.Confirm {
					// Use DryRunDefault from settings
					a.confirm = ui.NewConfirmModel(minutes, a.config.Settings.DryRunDefault)
					a.screen = ScreenConfirm
					return a, nil
				} else {
					// Skip confirmation and start immediately with DryRunDefault setting
					err := a.startShutdown(minutes, a.config.Settings.DryRunDefault)
					if err != nil {
						a.err = err.Error()
						return a, nil
					}
					// Go to active screen
					a.screen = ScreenActive
					a.active.Refresh(a.config)
					return a, a.active.Init()
				}
			}
		case "d":
			// Delete selected history item
			selected := a.history.GetSelectedHistory()
			if selected != nil {
				a.config.DeleteHistory(selected.ID)
				a.config.Save()
				a.history.Refresh(a.config)
				return a, nil
			}
		}
	}

	a.history, cmd = a.history.Update(msg)
	return a, cmd
}

// updateSettings handles updates for the settings screen
func (a *App) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Track previous language
	prevLang := a.config.Settings.Language

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Save settings and go back to home
			a.config.Save()
			a.screen = ScreenHome
			return a, nil
		}
	}

	a.settings, cmd = a.settings.Update(msg)

	// Check if language changed and reload translations
	if a.config.Settings.Language != prevLang {
		_ = i18n.SetLanguage(a.config.Settings.Language)
	}

	// Always save settings after update (for toggles)
	a.config.Save()

	return a, cmd
}

// startShutdown starts a shutdown timer
func (a *App) startShutdown(minutes int, dryRun bool) error {
	// Cancel any existing job first
	if a.config.ActiveJob != nil {
		_ = a.executor.Cancel(dryRun)
	}

	// Calculate job info
	jobInfo := shutdown.CalculateJobInfo(minutes)

	// Schedule shutdown
	command, err := a.executor.Schedule(minutes, dryRun)
	if err != nil {
		// Add to history as failed
		h := config.History{
			ID:              utils.GenerateID(),
			CreatedAt:       time.Now(),
			DurationSeconds: minutes * 60,
			ScheduledFor:    jobInfo.EndTime,
			Status:          config.StatusFailed,
			OS:              a.executor.GetOS(),
			Command:         command,
		}
		a.config.AddHistory(h)
		a.config.Save()
		return err
	}

	// Update config with active job
	a.config.ActiveJob = &config.ActiveJob{
		StartTime:   jobInfo.StartTime,
		EndTime:     jobInfo.EndTime,
		DurationSec: jobInfo.DurationSec,
		Command:     command,
	}

	// Add to history
	status := config.StatusOK
	if dryRun {
		status = config.StatusDryRun
	}
	h := config.History{
		ID:              utils.GenerateID(),
		CreatedAt:       jobInfo.StartTime,
		DurationSeconds: jobInfo.DurationSec,
		ScheduledFor:    jobInfo.EndTime,
		Status:          status,
		OS:              a.executor.GetOS(),
		Command:         command,
	}
	a.config.AddHistory(h)

	// Save config
	return a.config.Save()
}

// cancelShutdown cancels the current shutdown timer
func (a *App) cancelShutdown() error {
	if a.config.ActiveJob == nil {
		return nil
	}

	// Determine if it was a dry-run
	dryRun := false
	if len(a.config.History) > 0 && a.config.History[0].Status == config.StatusDryRun {
		dryRun = true
	}

	// Cancel the shutdown
	err := a.executor.Cancel(dryRun)
	if err != nil {
		// Update history status to failed
		a.config.UpdateHistoryStatus(config.StatusFailed)
	} else {
		// Update history status to cancelled
		a.config.UpdateHistoryStatus(config.StatusCancelled)
	}

	// Clear active job
	a.config.ActiveJob = nil

	// Save config
	saveErr := a.config.Save()
	if saveErr != nil {
		return saveErr
	}

	return err
}

// Cleanup performs cleanup tasks before application exit
func (a *App) Cleanup() {
	if a.config.ActiveJob != nil {
		// Check if job has expired
		now := time.Now()
		if now.After(a.config.ActiveJob.EndTime) {
			// Job has expired, clear it
			a.config.ActiveJob = nil
			_ = a.config.Save()
		}
	}
}
