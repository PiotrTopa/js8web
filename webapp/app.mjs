import ChatWindow from './chat-window.mjs'
import StatusBar from './status-bar.mjs'
import LoginPage from './login-page.mjs'
import ToastContainer from './toast-container.mjs'
import axios from 'axios'

const WS_RECONNECT_INTERVAL_MS = 3000

export default {
    components: {
        ChatWindow,
        StatusBar,
        LoginPage,
        ToastContainer,
    },
    data() {
        return {
            stationInfo: {},
            rigStatus: {},
            pttActive: false,
            wsConnected: false,
            ws: null,
            wsReconnectTimer: null,
            authenticated: false,
            authUser: null,
            authChecked: false,
        }
    },
    created() {
        this.checkAuth()
        this.fetchInitialData()
        this.connectWebSocket()
    },
    beforeUnmount() {
        if (this.wsReconnectTimer) {
            clearTimeout(this.wsReconnectTimer)
        }
        if (this.ws) {
            this.ws.close()
        }
    },
    methods: {
        checkAuth() {
            axios.get('/api/auth/check').then(response => {
                if (response.data.ok) {
                    this.authenticated = true
                    this.authUser = { username: response.data.username, role: response.data.role }
                }
            }).catch(() => {}).finally(() => {
                this.authChecked = true
            })
        },
        handleLogin(user) {
            this.authenticated = true
            this.authUser = user
            this.showToast({ type: 'success', message: 'Logged in as ' + user.username })
        },
        handleLogout() {
            axios.post('/api/auth/logout').then(() => {
                this.authenticated = false
                this.authUser = null
                this.showToast({ type: 'info', message: 'Logged out' })
            }).catch(() => {
                this.showToast({ type: 'error', message: 'Logout failed' })
            })
        },
        showToast(toast) {
            this.$refs.toasts?.add(toast)
        },
        fetchInitialData() {
            axios.get('/api/station-info').then(response => {
                this.stationInfo = response.data
            }).catch(() => {})
            axios.get('/api/rig-status').then(response => {
                this.rigStatus = response.data
            }).catch(() => {})
        },
        connectWebSocket() {
            if (this.ws) {
                this.ws.close()
                this.ws = null
            }

            const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
            const ws = new WebSocket(protocol + '//' + location.host + '/ws/events')

            ws.onopen = () => {
                this.wsConnected = true
                if (this.wsReconnectTimer) {
                    clearTimeout(this.wsReconnectTimer)
                    this.wsReconnectTimer = null
                }
                // Refresh cached data on reconnect
                this.fetchInitialData()
            }

            ws.onmessage = (message) => {
                const eventData = JSON.parse(message.data)

                // Update local state for status events
                if (eventData.EventType === 'event') {
                    if (eventData.WsType === 'STATION.INFO') {
                        this.stationInfo = eventData.Event
                    } else if (eventData.WsType === 'RIG.STATUS') {
                        this.rigStatus = eventData.Event
                    } else if (eventData.WsType === 'RIG.PTT') {
                        this.pttActive = eventData.Event.Enabled
                    }
                }

                // Broadcast to child components
                const event = new CustomEvent('event', { detail: eventData })
                window.dispatchEvent(event)
            }

            ws.onclose = () => {
                this.wsConnected = false
                this.scheduleReconnect()
            }

            ws.onerror = () => {
                // onclose will fire after onerror, reconnect handled there
            }

            this.ws = ws
        },
        scheduleReconnect() {
            if (this.wsReconnectTimer) return
            this.wsReconnectTimer = setTimeout(() => {
                this.wsReconnectTimer = null
                this.connectWebSocket()
            }, WS_RECONNECT_INTERVAL_MS)
        },
    },
    template: `
    <template v-if="!authChecked">
        <div class="d-flex align-items-center justify-content-center vh-100">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    </template>
    <template v-else-if="!authenticated">
        <LoginPage @login="handleLogin" />
    </template>
    <template v-else>
        <div class="d-flex flex-column vh-100">
            <StatusBar :stationInfo="stationInfo" :rigStatus="rigStatus" :connected="wsConnected" :authUser="authUser" @logout="handleLogout" />
            <div class="ptt-indicator" v-if="pttActive"><i class="bi bi-broadcast"></i> TX</div>
            <div class="flex-fill d-flex flex-column overflow-hidden p-2">
                <ChatWindow @toast="showToast" />
            </div>
        </div>
    </template>
    <ToastContainer ref="toasts" />
`
}