import Chat from './chat.mjs'

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
    <h3>Rig status</h3>
    <p>
        {{ rigStatus }}
    </p>

    <br />

    <h3>Station info</h3>
    <p>
        {{ stationInfo }}
    </p>

    <div>
    <Chat />
    </div>
`
}

