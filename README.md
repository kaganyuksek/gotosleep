# GoToSleep (gts)

Terminal-based shutdown timer application built with Go and Bubble Tea.

## Features

- 🎨 Beautiful terminal UI with intuitive navigation
- ⏱️ Quick presets (15m, 30m, 45m, 60m, 90m, 120m)
- ⌨️ Flexible duration input (90, 1h30m, 00:45, etc.)
- 📊 Real-time countdown with progress bar
- 📜 History tracking of all shutdown operations
- ⚙️ Configurable settings
- 🔒 Confirmation dialog with dry-run mode
- 🖥️ Cross-platform support (Windows, Linux, macOS)

## Installation

### Quick Install (Recommended)

Install directly from GitHub:

```bash
go install github.com/kaganyuksek/gotosleep/cmd/gts@latest
```

This will install `gts` to your `$GOPATH/bin` directory. Make sure this directory is in your PATH.

### Build from source

```bash
go build -o gts ./cmd/gts
```

### Windows

```bash
go build -o gts.exe ./cmd/gts
```

## Usage

After installation with `go install`, simply run:

```bash
gts
```

Or if you built from source:

```bash
./gts
```

Or on Windows:

```bash
gts.exe
```

### Navigation

**Home Screen:**

- `1-6`: Select quick preset
- Type duration and `Enter`: Start timer
- `h`: View history
- `s`: Open settings
- `a`: Go to active countdown (if running)
- `q`: Quit

**Confirm Dialog:**

- `Y`: Confirm and start shutdown
- `N` or `Esc`: Cancel
- `D`: Toggle dry-run mode

**Active Countdown:**

- `c`: Cancel shutdown
- `e`: Edit (cancel and create new timer)
- `h`: View history
- `Esc`: Return to home

**History Screen:**

- `↑↓`: Navigate list
- `Enter`: Restart selected timer
- `d`: Delete selected entry
- `Esc`: Go back

**Settings Screen:**

- `↑↓`: Navigate options
- `Space/Enter`: Toggle setting
- `Esc`: Save and go back

## Duration Formats

The application accepts various duration formats:

- `90` - 90 minutes
- `90m` - 90 minutes
- `1h30m` - 1 hour 30 minutes
- `2h` - 2 hours
- `00:45` - 45 minutes
- `1:20` - 1 hour 20 minutes

## Configuration

Configuration is stored in:

- **Windows:** `%AppData%\gts\state.json`
- **Linux:** `~/.config/gts/state.json`
- **macOS:** `~/Library/Application Support/gts/state.json`

## Platform Notes

### Windows

Works out of the box with no special permissions required.

### Linux/macOS

Shutdown commands require `sudo` privileges. You have two options:

1. Run the application with sudo:

   ```bash
   sudo gts
   ```

2. Configure sudoers to allow shutdown without password (advanced):
   ```bash
   # Add to /etc/sudoers.d/shutdown
   yourusername ALL=(ALL) NOPASSWD: /sbin/shutdown
   ```

### Dry-Run Mode

Enable dry-run mode to test the application without actually scheduling a shutdown. In this mode:

- No shutdown command is executed
- All UI features work normally
- History entries are marked as "dry-run"

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
