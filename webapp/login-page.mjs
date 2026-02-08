import axios from 'axios'

export default {
    emits: ['login'],
    data() {
        return {
            username: '',
            password: '',
            error: '',
            loading: false,
        }
    },
    mounted() {
        this.$refs.usernameInput?.focus()
    },
    methods: {
        submit() {
            if (!this.username || !this.password || this.loading) return

            this.loading = true
            this.error = ''

            axios.post('/api/auth/login', {
                username: this.username,
                password: this.password,
            }).then(response => {
                if (response.data.ok) {
                    this.$emit('login', {
                        username: response.data.username,
                        role: response.data.role,
                    })
                } else {
                    this.error = response.data.error || 'Login failed'
                }
            }).catch(err => {
                this.error = err.response?.data?.error || 'Connection error'
            }).finally(() => {
                this.loading = false
            })
        },
        handleKeydown(e) {
            if (e.key === 'Enter') {
                e.preventDefault()
                this.submit()
            }
        }
    },
    template: `
    <div class="login-overlay d-flex align-items-center justify-content-center vh-100">
        <div class="login-card card shadow">
            <div class="card-body p-4">
                <h4 class="card-title text-center mb-4">
                    <i class="bi bi-broadcast"></i> js8web
                </h4>
                <div class="alert alert-danger" v-if="error">
                    <i class="bi bi-exclamation-triangle"></i> {{ error }}
                </div>
                <div class="mb-3">
                    <label class="form-label" for="login-username">Username</label>
                    <input type="text" class="form-control" id="login-username" v-model="username" @keydown="handleKeydown" ref="usernameInput" :disabled="loading" autocomplete="username">
                </div>
                <div class="mb-3">
                    <label class="form-label" for="login-password">Password</label>
                    <input type="password" class="form-control" id="login-password" v-model="password" @keydown="handleKeydown" :disabled="loading" autocomplete="current-password">
                </div>
                <button class="btn btn-primary w-100" @click="submit" :disabled="!username || !password || loading">
                    <span v-if="loading"><i class="bi bi-hourglass-split"></i> Signing in...</span>
                    <span v-else><i class="bi bi-box-arrow-in-right"></i> Sign In</span>
                </button>
            </div>
        </div>
    </div>
    `
}
