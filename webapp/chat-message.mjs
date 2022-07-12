import ChatRxPacket from './chat-rx-packet.mjs'
import ChatRxMessage from './chat-rx-message.mjs'

export default {
    props: ['message'],
    components: {
        ChatRxPacket,
        ChatRxMessage,
    },
    methods: {
    },
    template: `
        <ChatRxPacket v-if="message.Type === 'RX.ACTIVITY'" :message=message />
        <ChatRxMessage v-if="message.Type === 'RX.DIRECTED'" :message=message />
    `
}