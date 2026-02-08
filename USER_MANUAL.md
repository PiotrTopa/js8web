# js8web — User Manual

## What is js8web?

js8web is a web-based monitor for [JS8Call](http://js8call.com/), an amateur radio digital communication program. It runs alongside JS8Call on your computer (or on any machine that can reach JS8Call over the network) and provides a browser-accessible dashboard showing received messages, station information, and rig status in real time.

**Key features:**
- Real-time chat-style display of all received JS8Call messages
- Filterable message history by callsign or frequency
- Tabbed interface for monitoring multiple conversations simultaneously
- Color-coded signal quality indicators (SNR, speed, time drift)
- Scrollable history with automatic loading of older messages
- Single binary — no additional software or runtime dependencies needed
- Access from any device with a web browser on your local network

---

## Requirements

- **JS8Call** — installed, configured, and running with the TCP API enabled
- **Go 1.18+** and **GCC** — needed only for building from source
- A modern web browser (Chrome, Firefox, Edge, Safari)

---

## Installation

### Building from Source

```bash
# Clone the repository
git clone https://github.com/PiotrTopa/js8web.git
cd js8web

# Build the binary
go build -o js8web .
```

This produces a single `js8web` executable with the web interface embedded inside.

### Pre-built Binary

If pre-built binaries are available from the releases page, download the one for your operating system and architecture. No installation is needed — just place it somewhere convenient.

---

## JS8Call Configuration

Before starting js8web, ensure that JS8Call's TCP API is enabled:

1. Open **JS8Call**
2. Go to **File → Settings → Reporting**
3. Check **Enable TCP Server API**
4. Note the **TCP Server Port** (default: `2442`)
5. If js8web runs on a different machine, set the listening address to `0.0.0.0` or the appropriate network interface

---

## Running js8web

```bash
./js8web
```

On the first run, js8web will:
1. Create a new SQLite database file (`js8web.db` in the current directory)
2. Set up the database schema automatically
3. Create a default admin user (`admin` / `admin`)
4. Begin attempting to connect to JS8Call at `localhost:2442`
5. Start the web server on port `8080`

### Command-Line Options

```bash
./js8web -port 9090 -js8call-addr 192.168.1.100:2442 -db /var/lib/js8web/data.db
```

| Flag | Description | Default |
|------|-------------|---------|
| `-js8call-addr` | JS8Call TCP API address (host:port) | `localhost:2442` |
| `-reconnect-interval` | Seconds between reconnection attempts | `5` |
| `-db` | Path to SQLite database file | `./js8web.db` |
| `-port` | HTTP server port | `8080` |
| `-log-level` | Log level: debug, info, warn, error | `info` |

All options can also be set via environment variables:

| Variable | Description |
|----------|-------------|
| `JS8WEB_JS8CALL_ADDR` | JS8Call TCP API address |
| `JS8WEB_RECONNECT_SEC` | Reconnect interval |
| `JS8WEB_DB_PATH` | Database file path |
| `JS8WEB_PORT` | HTTP server port |
| `JS8WEB_LOG_LEVEL` | Log level (debug/info/warn/error) |

CLI flags take precedence over environment variables.

Open your browser and navigate to:

```
http://localhost:8080
```

If accessing from another device on your network, use the IP address of the machine running js8web:

```
http://<your-ip-address>:8080
```

---

## Using the Web Interface

### Login

When you first open js8web in your browser, you will see a login page. Enter your credentials to access the dashboard.

The default account created on first run is:
- **Username:** `admin`
- **Password:** `admin`

> ⚠️ **Change the default password** after first login using the Admin panel (see [User Management](#user-management) below).

After logging in, your session is stored as a browser cookie and will persist for 24 hours. You can log out at any time using the logout button in the status bar.

### User Roles

js8web supports three user roles with different permission levels:

| Role | Permissions |
|------|-------------|
| **Admin** | Full access: view messages, send messages, manage users |
| **Operator** | View messages and send messages to JS8Call |
| **Monitor** | View messages only (read-only access) |

The default `admin` user has the Admin role. New users can be created with any role from the Admin panel.

### Status Bar

The top of the page shows a dark status bar with:

- **Connection indicator** — green wifi icon when connected to the server, blinking red wifi-off icon when disconnected
- **Call** — your station callsign
- **Grid** — your grid square
- **Dial** — current dial frequency in MHz
- **Offset** — audio offset in Hz
- **Speed** — JS8Call speed mode (color-coded)
- **Selected** — currently selected callsign in JS8Call
- **Info** — station info text
- **User** — your logged-in username and a logout button (right side)

The status bar updates in real time as JS8Call reports changes.

### PTT Indicator

When JS8Call is transmitting, a red **TX** banner appears below the status bar, pulsing to draw attention.

### Connection Auto-Reconnect

If the WebSocket connection to the server is lost (network issue, server restart), js8web will automatically attempt to reconnect every 3 seconds. The status bar connection indicator will turn red and blink while disconnected. When reconnected, cached data is refreshed automatically.

### Main Chat View

Below the status bar, you see a chat-style interface showing JS8Call messages:

- **Directed messages** (`RX.DIRECTED`) appear as chat bubbles with the sender's callsign, recipient, grid square, timestamp, and message content
- **Messages directed to you** (`RX.DIRECTED.ME`) appear with a green background to stand out
- **Raw activity packets** (`RX.ACTIVITY`) appear as compact lines with timestamp and decoded text (can be toggled on/off)
- **Transmitted frames** (`TX.FRAME`) appear with a red left border, showing when your station transmitted

### Message Indicators

Each message displays several indicators:

| Icon | Meaning |
|------|---------|
| 📡 Frequency (clickable) | The offset frequency of the signal. Click to open a new tab filtered to that frequency slot. |
| 🔵→🟡→🔴 SNR | Signal-to-noise ratio. Color ranges from blue (weak, -30 dB) through yellow (moderate) to red (strong, +20 dB). |
| ⏭ Speed | JS8Call speed mode: **S**low, **N**ormal, **F**ast, **T**urbo, **U**ltra |
| ⏱ Time Drift | Timing drift in milliseconds |

### Tabs and Filtering

The interface supports a tabbed view for monitoring specific conversations:

- **All messages** — default tab showing everything
- **Callsign filter** — click the 🔍 search icon next to any callsign to open a new tab showing only messages to/from that station
- **Frequency filter** — click any frequency indicator to open a tab filtered to that 50 Hz frequency slot
- **Close tabs** — click the ✕ button on any filter tab to close it

### Settings

Click the ⚙ gear icon in the tab bar to access settings:

- **Show raw packets** — toggle display of `RX.ACTIVITY` raw decoded packets alongside directed messages

### Sending Messages

When logged in, a message input field appears at the bottom of the chat view:

1. Type your message in the input field
2. Press **Enter** or click the **Send** button
3. The message is sent to JS8Call's `TX.SEND_MESSAGE` API, which queues it for transmission
4. A toast notification confirms the message was queued
5. Once transmitted, the TX frame will appear in the chat in real-time

> **Note:** Messages are sent exactly as typed. JS8Call handles the encoding and transmission. You can address specific stations using JS8Call's message format (e.g., `CALLSIGN MSG ...`).

> **Note:** Sending messages requires the **Operator** or **Admin** role. Monitor users see the chat but cannot send messages.

### Scrolling and History

- Messages are loaded in pages of up to 100 items (including both received packets and transmitted frames)
- **Scroll up** to load older messages automatically
- **Scroll down** to load newer messages
- When at the bottom of the list, new messages appear in real time
- At the top of the history, "(No more messages)" is displayed
- **TX frames** are now shown in historical scrolling alongside RX packets, sorted by timestamp

### User Management

Administrators can manage user accounts from the **Admin** tab (shield icon) in the tab bar. This tab is only visible to users with the Admin role.

From the Admin panel you can:

- **View all users** — see username, role, and bio for all accounts
- **Create new users** — click "New User" to add an account with a username, password, role, and optional bio
- **Change roles** — use the dropdown next to any user to change their role (admin/operator/monitor)
- **Reset passwords** — click the key icon to set a new password for any user
- **Delete users** — click the trash icon to remove an account (you cannot delete your own account)

---

## Connection Status

js8web automatically connects to JS8Call and will keep retrying if the connection is lost. Check the terminal output for connection status messages:

```
Connected to JS8call         — successful connection
Disconnected from JS8call    — connection lost, will retry
Connection to JS8call failed — cannot reach JS8Call, retrying in 5 seconds
```

When connected, the browser will show incoming messages in real time via WebSocket. If the page is refreshed, historical messages are loaded from the database and live streaming resumes.

---

## Database

All received packets, spots, and transmitted frames are stored in a SQLite database (`js8web.db`). This file is created automatically on first run.

- The database can be backed up by simply copying the file while js8web is stopped
- To start fresh, delete the database file and restart js8web
- The database can be queried directly using any SQLite client for custom analysis

### Tables

| Table | Contents |
|-------|----------|
| `RX_PACKET` | All received activity and directed messages |
| `RX_SPOT` | Spot reports from other stations |
| `TX_FRAME` | Your transmitted frames (with tone data) |
| `STATION_INFO` | Snapshots of your station configuration |
| `USERS` | User accounts (for future authentication) |

---

## Troubleshooting

### js8web cannot connect to JS8Call

- Ensure JS8Call is running and the TCP API is enabled (see [JS8Call Configuration](#js8call-configuration))
- Verify the TCP port matches (default `2442`)
- If running on different machines, ensure the firewall allows the connection
- js8web will keep retrying automatically — check terminal output for errors

### Web interface shows no messages

- Verify js8web is connected to JS8Call (check terminal output)
- Ensure JS8Call is receiving signals (check JS8Call's own display)
- Try refreshing the browser page
- Check that you can reach `http://localhost:8080/api/station-info` — it should return JSON

### Browser shows blank page

- Check the browser's developer console (F12) for JavaScript errors
- Ensure you're using a modern browser that supports ES modules and import maps
- Verify the URL is correct (`http://localhost:8080`)

### Database errors on startup

- If the database becomes corrupted, stop js8web, delete `js8web.db`, and restart
- Ensure the directory is writable by the user running js8web

---

## Security Considerations

> ⚠️ **js8web does not currently support HTTPS/TLS encryption.**

- js8web requires login to send messages — unauthenticated users cannot access the dashboard
- The default admin password is `admin` — **change it as soon as possible**
- Sessions are stored in server memory (not persisted) — a server restart will log out all users
- Session cookies are `HttpOnly` and `SameSite=Strict` to mitigate XSS and CSRF
- Do not expose js8web directly to the internet without additional protection (reverse proxy with TLS, VPN, firewall rules)
- Consider binding to `localhost` only if remote access is not needed
- All traffic including login credentials is unencrypted without HTTPS

---

## Limitations (Current Version)

- **No HTTPS** — all traffic is unencrypted; use a reverse proxy for TLS
- **In-memory sessions** — sessions do not survive server restarts
- **RX spots** — spot reports are stored but not yet displayed in the UI
