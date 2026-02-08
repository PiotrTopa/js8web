import axios from 'axios'

export default {
    emits: ['toast'],
    data() {
        return {
            users: [],
            loading: true,
            showCreateForm: false,
            newUser: { username: '', password: '', role: 'monitor', bio: '' },
            creating: false,
            editPasswordId: null,
            newPassword: '',
            changingPassword: false,
        }
    },
    created() {
        this.fetchUsers()
    },
    methods: {
        fetchUsers() {
            this.loading = true
            axios.get('/api/users').then(response => {
                if (response.data.ok) {
                    this.users = response.data.users || []
                }
            }).catch(() => {
                this.$emit('toast', { type: 'error', message: 'Failed to load users' })
            }).finally(() => {
                this.loading = false
            })
        },
        createUser() {
            if (!this.newUser.username.trim() || !this.newUser.password.trim()) return
            this.creating = true
            axios.post('/api/users', this.newUser).then(response => {
                if (response.data.ok) {
                    this.$emit('toast', { type: 'success', message: 'User created' })
                    this.newUser = { username: '', password: '', role: 'monitor', bio: '' }
                    this.showCreateForm = false
                    this.fetchUsers()
                }
            }).catch(err => {
                const msg = err.response?.data?.error || 'Failed to create user'
                this.$emit('toast', { type: 'error', message: msg })
            }).finally(() => {
                this.creating = false
            })
        },
        updateRole(user, newRole) {
            axios.put('/api/users/' + user.id, { role: newRole, bio: user.bio }).then(response => {
                if (response.data.ok) {
                    user.role = newRole
                    this.$emit('toast', { type: 'success', message: 'Role updated' })
                }
            }).catch(err => {
                const msg = err.response?.data?.error || 'Failed to update role'
                this.$emit('toast', { type: 'error', message: msg })
            })
        },
        startChangePassword(userId) {
            this.editPasswordId = userId
            this.newPassword = ''
        },
        cancelChangePassword() {
            this.editPasswordId = null
            this.newPassword = ''
        },
        changePassword(userId) {
            if (!this.newPassword.trim()) return
            this.changingPassword = true
            axios.put('/api/users/' + userId + '/password', { password: this.newPassword }).then(response => {
                if (response.data.ok) {
                    this.$emit('toast', { type: 'success', message: 'Password updated' })
                    this.editPasswordId = null
                    this.newPassword = ''
                }
            }).catch(err => {
                const msg = err.response?.data?.error || 'Failed to change password'
                this.$emit('toast', { type: 'error', message: msg })
            }).finally(() => {
                this.changingPassword = false
            })
        },
        deleteUser(user) {
            if (!confirm('Delete user "' + user.name + '"? This cannot be undone.')) return
            axios.delete('/api/users/' + user.id).then(response => {
                if (response.data.ok) {
                    this.$emit('toast', { type: 'success', message: 'User deleted' })
                    this.fetchUsers()
                }
            }).catch(err => {
                const msg = err.response?.data?.error || 'Failed to delete user'
                this.$emit('toast', { type: 'error', message: msg })
            })
        },
        roleBadgeClass(role) {
            switch (role) {
                case 'admin': return 'bg-danger'
                case 'operator': return 'bg-primary'
                case 'monitor': return 'bg-secondary'
                default: return 'bg-dark'
            }
        }
    },
    template: `
    <div class="admin-panel p-3">
        <div class="d-flex align-items-center mb-3">
            <h5 class="mb-0"><i class="bi bi-people"></i> User Management</h5>
            <button class="btn btn-sm btn-success ms-auto" @click="showCreateForm = !showCreateForm">
                <i class="bi" :class="showCreateForm ? 'bi-x-lg' : 'bi-person-plus'"></i>
                {{ showCreateForm ? 'Cancel' : 'New User' }}
            </button>
        </div>

        <div v-if="showCreateForm" class="card mb-3">
            <div class="card-body">
                <h6 class="card-title">Create User</h6>
                <div class="row g-2 align-items-end">
                    <div class="col-sm-3">
                        <label class="form-label small">Username</label>
                        <input type="text" class="form-control form-control-sm" v-model="newUser.username" placeholder="Username">
                    </div>
                    <div class="col-sm-3">
                        <label class="form-label small">Password</label>
                        <input type="password" class="form-control form-control-sm" v-model="newUser.password" placeholder="Password">
                    </div>
                    <div class="col-sm-2">
                        <label class="form-label small">Role</label>
                        <select class="form-select form-select-sm" v-model="newUser.role">
                            <option value="monitor">Monitor</option>
                            <option value="operator">Operator</option>
                            <option value="admin">Admin</option>
                        </select>
                    </div>
                    <div class="col-sm-2">
                        <label class="form-label small">Bio</label>
                        <input type="text" class="form-control form-control-sm" v-model="newUser.bio" placeholder="Optional">
                    </div>
                    <div class="col-sm-2">
                        <button class="btn btn-sm btn-primary w-100" @click="createUser" :disabled="creating || !newUser.username.trim() || !newUser.password.trim()">
                            <i class="bi" :class="creating ? 'bi-hourglass-split' : 'bi-check-lg'"></i> Create
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="loading" class="text-center py-3">
            <div class="spinner-border spinner-border-sm" role="status"></div> Loading...
        </div>

        <table v-else class="table table-sm table-hover mb-0">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Username</th>
                    <th>Role</th>
                    <th>Bio</th>
                    <th class="text-end">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="user in users" :key="user.id">
                    <td class="align-middle">{{ user.id }}</td>
                    <td class="align-middle fw-bold">{{ user.name }}</td>
                    <td class="align-middle">
                        <select class="form-select form-select-sm d-inline-block w-auto" :value="user.role" @change="updateRole(user, $event.target.value)">
                            <option value="monitor">Monitor</option>
                            <option value="operator">Operator</option>
                            <option value="admin">Admin</option>
                        </select>
                    </td>
                    <td class="align-middle text-muted small">{{ user.bio }}</td>
                    <td class="align-middle text-end">
                        <template v-if="editPasswordId === user.id">
                            <div class="input-group input-group-sm d-inline-flex w-auto">
                                <input type="password" class="form-control form-control-sm" v-model="newPassword" placeholder="New password" style="max-width:140px" @keydown.enter="changePassword(user.id)" @keydown.escape="cancelChangePassword">
                                <button class="btn btn-sm btn-success" @click="changePassword(user.id)" :disabled="changingPassword || !newPassword.trim()"><i class="bi bi-check-lg"></i></button>
                                <button class="btn btn-sm btn-outline-secondary" @click="cancelChangePassword"><i class="bi bi-x-lg"></i></button>
                            </div>
                        </template>
                        <template v-else>
                            <button class="btn btn-sm btn-outline-secondary me-1" @click="startChangePassword(user.id)" title="Change password">
                                <i class="bi bi-key"></i>
                            </button>
                            <button class="btn btn-sm btn-outline-danger" @click="deleteUser(user)" title="Delete user" :disabled="user.name === $root.authUser?.username">
                                <i class="bi bi-trash"></i>
                            </button>
                        </template>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
}
