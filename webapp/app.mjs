import Chat from './chat.mjs'
import { createApp } from 'vue'

export default {
    components: {
        Chat
    },
    data() {
        return {
            stationInfo: {},
            rigStatus: {},
            rxPackets: {},
        }
    },
    created() {
        this.$nextTick(this.connectToWebsocketEvents())
    },
    methods: {
        fetchData() {
            axios.get('/api/station-info').then(response => {
                this.stationInfo = response.data;
            });
            axios.get('/api/rig-status').then(response => {
                this.rigStatus = response.data;
            });
        },
        connectToWebsocketEvents() {
            var ws = new WebSocket("ws://" + location.host + "/ws/events");
            ws.onmessage = message => {
                const eventData = JSON.parse(message.data);
                const event = new CustomEvent("event", { detail: eventData });
                window.dispatchEvent(event);
            }
        },
    },
    template: `
    <div>
        <Chat />
    </div>
`
}

