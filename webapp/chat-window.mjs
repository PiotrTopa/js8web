import Chat from './chat.mjs'

function uidGenerator() {
    var S4 = function () {
        return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
    };
    return (S4() + S4() + "-" + S4() + "-" + S4() + "-" + S4() + "-" + S4() + S4() + S4());
}

export default {
    props: ['message'],
    components: {
        Chat,
    },
    data() {
        return {
            activeTab: 'all',
            chats: [
                {
                    id: 'all',
                    label: 'All messages',
                    filter: {},
                },
                {
                    id: 'test',
                    label: 'Test tab',
                    filter: { callsign: '161rx02' },
                }
            ],
            uid: uidGenerator(),
            settingsShowRawPackets: true,
        }
        
    },
    methods: {
        activateTab(selected) {
            this.activeTab = selected
        },
        closeTab(id) {
            this.chats = this.chats.filter(e => e.id != id)
        },
    },
    template: `
    <ul class="nav nav-tabs">
        <template v-for="chat in chats" :id="chat.id">
            <li class="nav-item" :class="{active: activeTab == chat.id}">
                <a class="nav-link" :class="{active: activeTab == chat.id}" @click="activateTab(chat.id)" href="#">
                    {{ chat.label }}
                    <a class="btn btn-light btn-sm" v-if="chat.id != 'all'" @click="closeTab(chat.id)"><i class="bi bi-x"></i></a>
                </a>
            </li>
        </template>
        <li class="nav-item" :class="{active: activeTab == 'settings'}">
            <a class="nav-link" :class="{active: activeTab == 'settings'}" @click="activateTab('settings')" href="#"><i class="bi bi-gear"></i></a>
        </li>
    </ul>
    <template v-for="chat in chats">
        <Chat v-show="activeTab == chat.id" :filter="chat.filter" :showRawPackets="this.settingsShowRawPackets" />
    </template>
    <div v-show="activeTab == 'settings'">
        <div class="row">
            <div class="col-12">
                <div class="form-check form-switch settings">
                    <input class="form-check-input" type="checkbox" role="switch" :id="this.uid+'-show-raw-packets'" v-model="this.settingsShowRawPackets">
                    <label class="form-check-label" :for="this.uid+'-show-raw-packets'">Show raw packets</label>
                </div>
            </div>
        </div>
    </div>
    `
}