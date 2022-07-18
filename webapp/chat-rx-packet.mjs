import ChatRxHeaderIcons from './chat-rx-header-icons.mjs'

export default {
    props: ['message'],
    emits: ['frequencySelected'],
    components: {
        ChatRxHeaderIcons
    },
    methods: {
    },
    template: `
        <li class="clearfix packet">
            <div class="header">
                <span class="time">{{ new Date(message.Timestamp).toLocaleTimeString() }}</span>
                <ChatRxHeaderIcons :message=message @frequencySelected="e => $emit('frequencySelected', e)" />
            </div>
            <br />
            <div class="content">{{ message.Text }}</div>
        </li>
    `
}

                