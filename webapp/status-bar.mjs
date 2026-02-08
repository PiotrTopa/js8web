export default {
    props: ['stationInfo', 'rigStatus', 'connected'],
    template: `
    <div class="status-bar d-flex align-items-center px-3 py-1">
        <div class="status-section me-4 d-flex align-items-center">
            <span class="connection-indicator me-2" :class="connected ? 'connected' : 'disconnected'" :title="connected ? 'Connected' : 'Disconnected'">
                <i class="bi" :class="connected ? 'bi-wifi' : 'bi-wifi-off'"></i>
            </span>
        </div>

        <div class="status-section me-4" v-if="stationInfo.Call">
            <span class="status-label">Call</span>
            <span class="status-value fw-bold">{{ stationInfo.Call }}</span>
        </div>

        <div class="status-section me-4" v-if="stationInfo.Grid">
            <span class="status-label">Grid</span>
            <span class="status-value">{{ stationInfo.Grid }}</span>
        </div>

        <div class="status-section me-4" v-if="rigStatus.Dial">
            <span class="status-label">Dial</span>
            <span class="status-value">{{ formatFrequency(rigStatus.Dial) }}</span>
        </div>

        <div class="status-section me-4" v-if="rigStatus.Offset">
            <span class="status-label">Offset</span>
            <span class="status-value">{{ rigStatus.Offset }} Hz</span>
        </div>

        <div class="status-section me-3" v-if="rigStatus.Speed">
            <span class="status-label">Speed</span>
            <span class="status-value speed-badge" :class="'speed-' + rigStatus.Speed">{{ rigStatus.Speed }}</span>
        </div>

        <div class="status-section me-3" v-if="rigStatus.Selected">
            <span class="status-label">Selected</span>
            <span class="status-value">{{ rigStatus.Selected }}</span>
        </div>

        <div class="status-section ms-auto" v-if="stationInfo.Info">
            <span class="status-label">Info</span>
            <span class="status-value text-muted small">{{ stationInfo.Info }}</span>
        </div>
    </div>
    `,
    methods: {
        formatFrequency(dialHz) {
            const mhz = dialHz / 1000000
            return mhz.toFixed(3) + ' MHz'
        }
    }
}
