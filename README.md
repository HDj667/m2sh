# m2sh
m2sh is a CLI mattermost shell. Personal learning/research go programming.

## Configuration

m2sh supports configuration through both INI files and environment variables.

### Configuration Files

m2sh looks for configuration files in the following locations (in order):
1. `./m2sh.ini` (current directory)
2. `~/.m2sh.ini` (home directory)
3. `~/.config/m2sh.ini` (XDG config directory)
4. `/etc/m2sh/m2sh.ini` (system-wide)

Example configuration file (`m2sh.ini`):
```ini
# Mattermost server URL (required)
url = https://mattermost.example.com

# Mattermost username (required)
username = your-username

# Mattermost password (optional, will be prompted if not set)
password = 
```

See `m2sh.ini.example` for a complete example.

### Environment Variables

Environment variables have **higher priority** than configuration file values:
- `MM_URL` - Mattermost server URL
- `MM_USERNAME` - Mattermost username
- `MM_PASSWORD` - Mattermost password (optional)

### Priority

Configuration values are loaded in the following order (highest priority first):
1. Environment variables (`MM_*`)
2. Configuration file values
3. Interactive prompts (for password if not set)
