import Chat from './chat.mjs'

var ws = new WebSocket("ws://" + location.host + "/ws/updates");
ws.onmessage = message => {
    var event = JSON.parse(message.data);
    console.log("ws", event)
}

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

