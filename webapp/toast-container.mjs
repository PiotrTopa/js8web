export default {
    data() {
        return {
            toasts: [],
            nextId: 1,
        }
    },
    methods: {
        add(toast) {
            const id = this.nextId++
            const type = toast.type || 'info'
            const duration = toast.duration || (type === 'error' ? 6000 : 3000)
            this.toasts.push({ id, type, message: toast.message })
            setTimeout(() => this.remove(id), duration)
        },
        remove(id) {
            this.toasts = this.toasts.filter(t => t.id !== id)
        },
        iconClass(type) {
            switch (type) {
                case 'success': return 'bi-check-circle-fill text-success'
                case 'error': return 'bi-exclamation-triangle-fill text-danger'
                case 'warning': return 'bi-exclamation-circle-fill text-warning'
                default: return 'bi-info-circle-fill text-info'
            }
        }
    },
    template: `
    <div class="toast-container">
        <transition-group name="toast">
            <div v-for="toast in toasts" :key="toast.id" class="toast-item" :class="'toast-' + toast.type" @click="remove(toast.id)">
                <i class="bi" :class="iconClass(toast.type)"></i>
                <span class="toast-message">{{ toast.message }}</span>
                <button class="btn-close btn-close-sm ms-2" @click.stop="remove(toast.id)"></button>
            </div>
        </transition-group>
    </div>
    `
}
