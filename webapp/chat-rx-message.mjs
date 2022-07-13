import ChatRxHeaderIcons from './chat-rx-header-icons.mjs'

export default {
    props: ['message'],
    components: {
        ChatRxHeaderIcons
    },
    methods: {
    },
    template: `
        <li class="clearfix message my-message">
            <div class="header">
                <span class="time">{{ new Date(message.Timestamp).toLocaleString() }}</span>
                <ChatRxHeaderIcons :message=message />
                <br />

                <span class="from">{{ message.From }}</span>
                <span class="grid" v-if=message.Grid><i class="bi bi-globe"></i>{{ message.Grid }}</span>
            </div>
            <div class="content">{{ message.Text }}</div>
        </li>
    `
}

                