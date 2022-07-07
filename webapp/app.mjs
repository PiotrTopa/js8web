import axios from 'axios'

export default {
    data() {
        return {
            stationInfo: {},
            rigStatus: {},
            days: {},
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
            axios.get('/api/rx-packets/list-days').then(response => {
                this.days = response.data;
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

    <br />

    <button @click="fetchData()">
        Update
    </button>`
}

