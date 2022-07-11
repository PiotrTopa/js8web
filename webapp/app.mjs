import axios from 'axios'
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

        fetchRxPackets() {
            const today = new Date()
            const yesterday = new Date()
            yesterday.setDate(today.getDate() - 1)

            axios.get('/api/rx-packets', {
                params: {
                    from: yesterday.toISOString(),
                    to: today.toISOString()
                }
            }).then(response => {
                this.rxPackets = response.data;
            })
        }
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

    <Chat :messages=this.rxPackets />

    <button @click="fetchData()">
        Update
    </button>
    <button @click="fetchRxPackets()">
        Fetch RxPackets
    </button>`
}

