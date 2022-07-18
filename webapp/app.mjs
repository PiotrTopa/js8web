import ChatWindow from './chat-window.mjs'

export default {
    components: {
        ChatWindow
    },
    data() {
        return {
            stationInfo: {},
            rigStatus: {},
            rxPackets: {}
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
    <div class="container-fluid main_container d-flex vh-100">
        <div class="row flex-fill">
            <div class="col-12 h-100">
                <div class="row h-75">
                    <div class="col-sm-12 col-md-10 d-flex flex-column mh-100">
                        <ChatWindow />
                    </div>
                    <div class="col-md-2">
                        <!-- Button for information -->
                        Info
                    </div>
                </div>
                <div class="row h-25">
                    <div class="col-sm-12">
                        <!-- Button for information -->
                        Info
                    </div>
                </div>
            </div>
        </div>
    </div>
`
}

