# js8web

Web-based monitor and control interface for [JS8Call](http://js8call.com/) — an amateur radio digital communication application.

js8web connects to a running JS8Call instance via its TCP API, captures received messages and station events in real time, stores them in a local SQLite database, and presents everything through a browser-accessible chat-style dashboard.

## Quick Start

```bash
# Build (requires Go 1.18+ and GCC)
go build -o js8web .

# Run (JS8Call must be running with TCP API enabled on port 2442)
./js8web

# Open in browser
# http://localhost:8080
```

## Features

- 📡 Real-time display of all received JS8Call messages and activity
- 🔍 Filter by callsign or frequency in dynamic tabs
- 📊 Color-coded SNR, speed mode, and time drift indicators
- 📜 Scrollable message history with automatic pagination
- 💾 SQLite database for persistent logging of all activity
- 📦 Single binary with embedded web interface — no dependencies at runtime

## Documentation

- **[User Manual](USER_MANUAL.md)** — installation, configuration, and usage guide
- **[Development Documentation](DEVELOPMENT.md)** — architecture, code reference, API docs, and implementation status

## License

See [LICENSE](LICENSE) for details.
