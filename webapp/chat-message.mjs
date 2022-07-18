import ChatRxPacket from './chat-rx-packet.mjs'
import ChatRxMessage from './chat-rx-message.mjs'

export default {
    props: ['message', 'showRawPackets'],
    emits: ['callsignSelected', 'frequencySelected'],
    components: {
        ChatRxPacket,
        ChatRxMessage,
    },
    methods: {
    },
    template: `
        <ChatRxPacket v-if="showRawPackets && message.Type === 'RX.ACTIVITY'" :message=message @frequencySelected="e => $emit('frequencySelected', e)" />
        <ChatRxMessage @callsignSelected="e => $emit('callsignSelected', e)" @frequencySelected="e => $emit('frequencySelected', e)" v-if="message.Type === 'RX.DIRECTED'" :message=message />
    `
}