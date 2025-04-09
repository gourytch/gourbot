# Feature Requests

## Implemented Features

### Storage Module
- SQLite-based storage implemented.
- Methods for opening, closing, and deleting the database.
- Method to create necessary tables in the database.
- Method to log Telegram API interactions (`AddTgRecord`).

### Configuration Module
- Added a `DbPath` field to the configuration for specifying the database path.
- Common `defaultPrefix` used for log and database paths.

### Testing
- Unit tests implemented for the `storage` module.
- In-memory SQLite used for testing.

### Telegram Bot
- Dynamic command registration implemented.
- `/list` command displays available commands dynamically.
- `/stop` command with proper shutdown handling.

### Logging
- Improved logging for better debugging and monitoring.

## Pending Feature Requests

### Telegram Bot
- Add more commands for advanced bot functionality.

### Documentation
- Keep the feature request list updated.