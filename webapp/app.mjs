import Chat from './chat.mjs'

var ws = new WebSocket("ws://" + location.host + "/ws/updates");
ws.onmessage = event => console.log("ws-message", event)

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
    methods: {
        fetchData() {
            axios.get('/api/station-info').then(response => {
                this.stationInfo = response.data;
            });
            axios.get('/api/rig-status').then(response => {
                this.rigStatus = response.data;
            });
        },
    },
    template: `
    <div>
        <Chat />
    </div>
`
}

