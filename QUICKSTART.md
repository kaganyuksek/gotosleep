# Quick Start Guide

## First Time Setup

1. **Build the application:**

   ```bash
   go build -o sdt.exe ./cmd/sdt
   ```

2. **Run the application:**
   ```bash
   ./sdt.exe
   ```

## Basic Usage Flow

### Setting a Quick Timer

1. Launch `gts.exe`
2. Press `1` through `6` to select a preset (15m, 30m, 45m, 60m, 90m, 120m)
3. Press `Enter`
4. Press `Y` to confirm
5. Watch the countdown!

### Setting a Custom Timer

1. Launch `gts.exe`
2. Type your duration (examples: `90`, `1h30m`, `00:45`)
3. Press `Enter`
4. Press `Y` to confirm
5. Watch the countdown!

### Using Dry-Run Mode (Testing)

1. When on the confirm screen, press `D` to toggle dry-run mode
2. Press `Y` to confirm
3. The UI will work normally but no shutdown will be scheduled
4. Perfect for testing!

### Viewing History

1. From the home screen, press `h`
2. Use `↑↓` to navigate through past timers
3. Press `Enter` to restart a previous timer
4. Press `d` to delete an entry
5. Press `Esc` to go back

### Canceling a Shutdown

1. Press `c` from the active countdown screen
2. Or press `e` to cancel and immediately create a new timer

### Adjusting Settings

1. From the home screen, press `s`
2. Use `↑↓` to navigate
3. Press `Space` or `Enter` to toggle settings
4. Press `Esc` to save and go back

## Tips & Tricks

- **Quick Access:** If a shutdown is running, press `a` from home to jump to the countdown
- **Flexible Input:** You can type durations in many formats: `90`, `90m`, `1h30m`, `1:30`
- **History Reuse:** Use history to quickly restart commonly-used timers
- **Safe Testing:** Always enable dry-run mode when testing to avoid accidental shutdowns

## Keyboard Shortcuts Cheat Sheet

### Global

- `q` or `Ctrl+C`: Quit application

### Home Screen

- `1-6`: Quick presets
- `Enter`: Start timer
- `h`: History
- `s`: Settings
- `a`: Active countdown

### Confirm Dialog

- `Y`: Yes/Confirm
- `N` or `Esc`: No/Cancel
- `D`: Toggle dry-run

### Active Countdown

- `c`: Cancel
- `e`: Edit (cancel & new)
- `h`: History
- `Esc`: Back to home

### History

- `↑↓` or `j/k`: Navigate
- `Enter`: Restart timer
- `d`: Delete entry
- `Esc`: Back

### Settings

- `↑↓` or `j/k`: Navigate
- `Space/Enter`: Toggle
- `Esc`: Save & back

## Platform-Specific Notes

### Windows ✅

Works perfectly without any special setup!

### Linux/macOS ⚠️

Requires sudo for shutdown commands:

```bash
sudo ./sdt
```

Or configure sudoers (advanced users):

```bash
echo "$USER ALL=(ALL) NOPASSWD: /sbin/shutdown" | sudo tee /etc/sudoers.d/shutdown
```

## Troubleshooting

**Problem:** "Permission denied" on Linux/macOS
**Solution:** Run with `sudo` or configure sudoers

**Problem:** Want to test without actual shutdown
**Solution:** Enable dry-run mode with `D` key in confirm dialog

**Problem:** Countdown not visible
**Solution:** Resize terminal to at least 80x24 characters

**Problem:** Application won't start
**Solution:** Make sure Go 1.21+ is installed and dependencies are downloaded

## Configuration File

Configuration is automatically saved at:

- Windows: `%AppData%\sdt\state.json`
- Linux: `~/.config/sdt/state.json`
- macOS: `~/Library/Application Support/sdt/state.json`

You can manually edit this file if needed, but be careful with the JSON format!

## Need Help?

Check the main README.md for more detailed documentation.
