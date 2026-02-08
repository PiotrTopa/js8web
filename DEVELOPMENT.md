# js8web — Development Documentation

## Project Overview

**js8web** is a web-based monitor and control interface for [JS8Call](http://js8call.com/), an amateur radio digital communication application built on top of the FT8 protocol. js8web connects to a running JS8Call instance via its TCP API, captures all incoming events (received packets, spots, rig status changes, etc.), persists them to a local SQLite database, and exposes them through a real-time web UI and REST/WebSocket API.

The goal is to provide a remote, browser-accessible dashboard for monitoring and eventually controlling a JS8Call station — useful for headless operation, multi-operator setups, and logging/archival.

---

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Backend language | Go 1.18 |
| Web framework | `net/http` standard library + `labstack/echo` (imported but unused currently; `net/http` `ServeMux` is the active router) |
| WebSocket | `gorilla/websocket` |
| Database | SQLite 3 (`mattn/go-sqlite3`) |
| Logging | `go.uber.org/zap` |
| Frontend framework | Vue.js 3 (ESM via CDN, no build step) |
| CSS framework | Bootstrap 5.2 + Bootstrap Icons |
| HTTP client (frontend) | axios (CDN import map) |
| Static file serving | Go `embed` (webapp directory is embedded into the binary) |

---

## Architecture

```
┌─────────────┐       TCP JSON        ┌──────────────────────────────────────────┐
│   JS8Call    │ ◄──────────────────── │              js8web (Go binary)          │
│  Application │ ─────────────────────►│                                          │
└─────────────┘   localhost:2442       │  ┌──────────────────────────────────┐    │
                                       │  │  js8call.go                      │    │
                                       │  │  TCP client with auto-reconnect  │    │
                                       │  └────────┬───────────┬─────────────┘    │
                                       │           │ incoming  │ outgoing         │
                                       │           ▼           │                  │
                                       │  ┌────────────────┐   │                  │
                                       │  │ dispatcher.go   │   │                  │
                                       │  │ Event routing   │   │                  │
                                       │  └───┬────────┬───┘   │                  │
                                       │      │        │        │                  │
                                       │      ▼        ▼        │                  │
                                       │  ┌──────┐ ┌───────┐   │                  │
                                       │  │SQLite│ │  WS   │   │                  │
                                       │  │  DB  │ │Events │   │                  │
                                       │  └──────┘ └───┬───┘   │                  │
                                       │               │        │                  │
                                       │  ┌────────────▼────────▼─────────────┐   │
                                       │  │     webappServer.go               │   │
                                       │  │  HTTP server (:8080)              │   │
                                       │  │  ┌─────────┐ ┌──────┐ ┌───────┐  │   │
                                       │  │  │REST API  │ │  WS  │ │Static │  │   │
                                       │  │  │/api/*    │ │/ws/* │ │webapp │  │   │
                                       │  │  └─────────┘ └──────┘ └───────┘  │   │
                                       │  └──────────────────────────────────┘   │
                                       └──────────────────────────────────────────┘
                                                         │
                                                         ▼
                                               ┌─────────────────┐
                                               │  Browser (Vue3) │
                                               │  Chat-style UI  │
                                               └─────────────────┘
```

---

## Source File Reference

### Backend (Go, root package `main`)

| File | Responsibility |
|------|---------------|
| `main.go` | Entry point. Initializes logger, config, DB, channels, JS8Call connection, dispatcher, WebSocket session container, and HTTP server. Handles graceful shutdown via OS signals (SIGINT/SIGTERM). |
| `const.go` | Configuration: CLI flags, environment variable parsing, defaults. Embeds `res/initDb.sql` via `//go:embed`. |
| `js8call.go` | Manages the persistent TCP connection to JS8Call. Auto-reconnects on failure. Reads newline-delimited JSON events from JS8Call and writes outgoing events. |
| `dispatcher.go` | Central event router. Applies a fix for the ambiguous `STATION.STATUS` event. Dispatches each event type to its specific notifier function, producing `WebsocketEvent` or `DbObj` items. |
| `rxActivity.go` | Notifier for `RX.ACTIVITY`, `RX.DIRECTED`, `RX.DIRECTED.ME` events → creates `RxPacketObj` for DB. Also handles `RX.SPOT` → creates `RxSpotObj`. |
| `rigStatus.go` | Notifier for `RIG.STATUS` (synthesized type) → updates in-memory `rigStatusCache` and emits WS event on change. Also handles `RIG.PTT`. |
| `stationInfo.go` | Notifier for `STATION.CALLSIGN`, `STATION.GRID`, `STATION.INFO`, `STATION.STATUS` → updates in-memory `stationInfoCache`, emits WS event and persists to DB. Initializes cache from last DB record on startup. |
| `txActivity.go` | Notifier for `TX.FRAME` → creates `TxFrameObj`, applies current rig status, saves to DB. |
| `db.go` | SQLite database initialization. Creates file and runs `initDb.sql` if DB does not exist. Creates default admin user. |
| `webappServer.go` | HTTP server setup. Registers REST API routes, WebSocket endpoint, and static file handler for the embedded webapp. Implements method routing. |
| `api.go` | REST API handler functions: `GET /api/station-info`, `GET /api/rig-status`, `GET /api/rx-packets`. |
| `websocket.go` | WebSocket upgrade handler and session management. Each connected browser gets a session; all sessions receive broadcast messages. |

### Model Package (`model/`)

| File | Responsibility |
|------|---------------|
| `js8callEvent.go` | Defines `Js8callEvent` and `Js8callEventParams` structs (JSON mapping of JS8Call TCP API). Declares all event type constants and WS type constants. Helper functions for channel calculation and speed naming. |
| `db.go` | Defines `DbObj` interface: objects that can be `Save()`d to DB and have a `WsType()`. |
| `websocketEvent.go` | Defines `WebsocketEvent` interface (anything with `WsType()`). |
| `websocketMessage.go` | Defines `WebsocketMessage` struct sent over WebSocket to browsers. |
| `rxPacket.go` | `RxPacketObj` — model for received packets. Insert, Scan, query logic. Supports filtered listing with pagination by timestamp (before/after). |
| `rxSpot.go` | `RxSpotObj` — model for RX spot reports. Insert logic. Stub for listing by days. |
| `txFrame.go` | `TxFrameObj` — model for transmitted frames. Stores tone data as JSON. Applies rig status before saving. |
| `rigStatus.go` | `RigStatusWsEvent` — in-memory rig status (dial freq, offset, speed, selected callsign). Not persisted. |
| `rigPtt.go` | `RigPttWsEvent` — PTT on/off event for WebSocket broadcast. |
| `stationInfo.go` | `StationInfoObj` / `StationInfoWsEvent` — station callsign, grid, info, status. Persisted with a "latest" flag pattern. |
| `user.go` | `User` model with SHA-256 password hashing. Default admin/admin user. Roles: admin, monitor, operator. Insert only — no login/auth flow implemented yet. |
| `utils.go` | Time conversion helpers: JS8Call millisecond timestamps ↔ `time.Time` ↔ SQLite RFC3339 strings. |

### Database Schema (`res/initDb.sql`)

| Table | Purpose |
|-------|---------|
| `USERS` | User accounts (name, password hash, role, bio) |
| `RX_PACKET` | Every received packet: timestamp, type, frequency info, SNR, speed, grid, from/to callsigns, text content, command/extra fields |
| `RX_SPOT` | Spot reports: call, grid, SNR, frequency info |
| `TX_FRAME` | Transmitted frames: frequency info, mode, speed, selected callsign, tone data |
| `STATION_INFO` | Station metadata snapshots: callsign, grid, info, status. Uses `LATEST=1` flag for current. |

All timestamp-bearing tables have indexes on `TIMESTAMP`.

### Frontend (`webapp/`)

| File | Responsibility |
|------|---------------|
| `index.html` | Single-page app shell. Loads Bootstrap, Vue 3, axios via CDN. Mounts Vue app. |
| `app.mjs` | Root Vue component. Fetches initial station info and rig status via REST API. Opens WebSocket with auto-reconnect logic (3s interval). Updates local state for station info, rig status, and PTT from WebSocket events. Dispatches events to child components via browser `CustomEvent`s. |
| `status-bar.mjs` | Status bar component showing connection state (wi-fi icon), station callsign, grid, dial frequency, offset, speed mode, selected callsign, and station info. |
| `chat-window.mjs` | Tab management component. Default "All messages" tab + dynamic filter tabs (by callsign or frequency). Settings tab with "show raw packets" toggle. |
| `chat.mjs` | Core chat/message list component. Infinite scroll (loads older/newer pages). Listens for `RX.PACKET` WebSocket events and appends new messages in real-time. Applies client-side filtering. |
| `chat-message.mjs` | Router component: renders `ChatRxPacket` for raw `RX.ACTIVITY`, `ChatRxMessage` for `RX.DIRECTED`/`RX.DIRECTED.ME` messages, and `ChatTxFrame` for transmitted frames. |
| `chat-rx-message.mjs` | Renders a directed message with sender callsign, recipient, grid, timestamp, SNR/speed/drift gauges, and message text. Messages directed to own station are visually highlighted. |
| `chat-rx-packet.mjs` | Renders a raw activity packet with timestamp and gauges. |
| `chat-tx-frame.mjs` | Renders a transmitted frame indicator with timestamp, frequency, speed, and selected callsign. |
| `chat-rx-header-icons.mjs` | Reusable gauge icons: frequency (clickable to filter), SNR (color-coded blue→yellow→red), speed indicator, time drift. |
| `style.css` | Chat-style layout. Flex-based full-height UI. Message bubbles, gauge styling, speed-color classes. |

---

## Data Flow

### Incoming (JS8Call → Browser)

1. **TCP Read** — `js8call.go` reads newline-delimited JSON from JS8Call TCP socket.
2. **Parse** — JSON is unmarshalled into `model.Js8callEvent`.
3. **Fix** — `dispatcher.go` renames ambiguous `STATION.STATUS` events that carry frequency info to `RIG.STATUS`.
4. **Dispatch** — Event is routed to the appropriate notifier based on type.
5. **Process** — Notifier creates either:
   - A `DbObj` (saved to SQLite, then broadcast to WS as `"object"` type), or
   - A `WebsocketEvent` (broadcast to WS as `"event"` type, not persisted), or both.
6. **Broadcast** — `mainDispatcher` sends `WebsocketMessage` to all connected WS sessions.
7. **Display** — Browser receives JSON via WebSocket, fires a `CustomEvent`, and Vue components update reactively.

### Outgoing (Browser → JS8Call)

The `outgoingEvents` channel exists in `main.go` and is wired through `js8call.go`, but **no outgoing commands are sent yet**. The chat input field exists in the HTML template but is hidden (`display: none`).

### REST API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `GET /api/station-info` | GET | Returns current station info (callsign, grid, info, status) from in-memory cache. |
| `GET /api/rig-status` | GET | Returns current rig status (dial, freq, offset, channel, speed, selected) from in-memory cache. |
| `GET /api/rx-packets` | GET | Returns up to 100 RX packets. Params: `startTime` (RFC3339), `direction` (`before`/`after`), optional `filter` (JSON with `Callsign` and/or `Freq.From`/`Freq.To`). |

### WebSocket

| Endpoint | Direction | Description |
|----------|-----------|-------------|
| `ws://host:8080/ws/events` | Server → Client | Broadcasts all state changes. Message format: `{ EventType, WsType, Event }`. EventType is `"object"` (persisted) or `"event"` (transient). |

---

## Configuration

Configuration is handled via CLI flags and environment variables (defined in `const.go`).
CLI flags take precedence over environment variables, which take precedence over defaults.

| CLI Flag | Environment Variable | Default | Description |
|----------|---------------------|---------|-------------|
| `-js8call-addr` | `JS8WEB_JS8CALL_ADDR` | `localhost:2442` | JS8Call TCP API address |
| `-reconnect-interval` | `JS8WEB_RECONNECT_SEC` | `5` | Seconds between reconnection attempts |
| `-db` | `JS8WEB_DB_PATH` | `./js8web.db` | SQLite database file path |
| `-port` | `JS8WEB_PORT` | `8080` | HTTP server listen port |

Run `./js8web -help` to see all options.

---

## Build & Run

```bash
# Prerequisites: Go 1.18+, GCC (for CGo/SQLite)
go build -o js8web .
./js8web
```

The webapp is embedded in the binary — no separate deployment needed.

---

## Current Implementation Status

### ✅ Working

- TCP connection to JS8Call with auto-reconnect
- Parsing of all major JS8Call event types
- SQLite persistence for RX packets, RX spots, TX frames, station info
- In-memory caching of rig status and station info
- REST API for station info, rig status, and paginated RX packet listing with filters
- WebSocket broadcast of real-time events to all connected browsers
- Vue 3 SPA with chat-style message display
- **Status bar** showing callsign, grid, dial frequency, offset, speed, selected callsign, and station info
- **Connection status indicator** (green wifi icon when connected, blinking red when disconnected)
- **WebSocket auto-reconnect** (3-second interval) with automatic data refresh on reconnect
- **PTT indicator** — red banner when transmitting
- **TX frame display** — transmitted frames shown in real-time in the chat
- **RX.DIRECTED.ME highlighting** — messages directed to own station visually distinguished (green background)
- Tab system with dynamic filter tabs (by callsign or frequency)
- Infinite scroll for message history
- Color-coded SNR, speed indicators, time drift gauges
- Raw packet toggle in settings
- Embedded static files (single binary deployment)
- **Configuration via CLI flags and environment variables** (`-port`, `-db`, `-js8call-addr`, `-reconnect-interval`)
- **Graceful shutdown** via SIGINT/SIGTERM signals
- **Thread-safe WebSocket session management** with `sync.RWMutex`
- **Non-blocking WebSocket broadcast** (slow clients don't block others)
- Proper JSON `Content-Type` headers on all API responses
- Proper `<!DOCTYPE html>` and viewport meta for mobile support

### 🚧 Partially Implemented / Stubbed

- **User model** — schema and model exist, default admin user is created, but **no authentication or authorization** is implemented. No login endpoint, no session management, no protected routes.
- **Outgoing events** — channel infrastructure exists end-to-end, but no UI or API to send commands to JS8Call.
- **RX Spot listing** — `RxSpotListDays()` function body is empty (stub).

### ❌ Not Yet Implemented

- Authentication & authorization (login, sessions, role-based access)
- Sending messages / commands to JS8Call from the web UI
- RX spot display / spot map
- TX frame historical loading (TX frames appear in real-time but not when scrolling through history)
- HTTPS / TLS support
- Mobile-responsive layout refinements
- Logging level configuration
- Unit / integration tests
- CI/CD pipeline
- Docker container / systemd service file

---

## Code Conventions

- **Go**: Standard Go project layout. Single `main` package with domain logic split by feature file. Model types in `model/` sub-package.
- **Frontend**: Vue 3 Composition-style components using `.mjs` ES modules loaded directly via browser import maps (no bundler). Component templates are inline template strings.
- **Database**: SQL queries are package-level `var` string constants. Prepared statements are used for all queries.
- **Naming**: JS8Call event types use dot notation (`RX.ACTIVITY`). Go structs use `Obj` suffix for DB-persisted models and `WsEvent` suffix for WebSocket-only events.
