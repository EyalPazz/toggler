# Toggler

A command-line interface for Toggl time tracking.

## Installation

1. Build the binary:
   ```bash
   go build -o toggler
   ```

## Configuration

### Option 1: Environment Variable
```bash
export TOGGLER_API_TOKEN="your_toggl_api_token_here"
```

### Option 2: Configuration File
Create a `toggler.yaml` file in your home directory (`~/.config/toggler.yaml`) or current directory:
```yaml
api_token: "your_toggl_api_token_here"
```

## Getting Your API Token

1. Go to https://toggl.com/app/profile
2. Scroll down to find your API token
3. Copy the token and use it in your configuration

## Usage

### Start a timer
```bash
./toggler start "Working on project"
./toggler start -d "Meeting with client"
```

### Check current timer
```bash
./toggler current
```

### Stop current timer
```bash
./toggler stop
```

### List recent entries
```bash
./toggler list              # Today's entries
./toggler list --days 7     # Last 7 days
```

### Help
```bash
./toggler --help
./toggler [command] --help
```

## Commands

- `start [description]` - Start a new time entry
- `stop` - Stop the currently running timer
- `current` - Show currently running timer info  
- `list --days N` - List recent time entries
- `help` - Show help information

## Troubleshooting

If you get API errors:
1. Verify your API token is correct
2. Check your internet connection
3. Ensure you have an active Toggl account

The CLI uses Toggl API v8. For more information, see https://engineering.toggl.com/docs/