import ChatRxPacket from './chat-rx-packet.mjs'
import ChatRxMessage from './chat-rx-message.mjs'
import ChatTxFrame from './chat-tx-frame.mjs'

export default {
    props: ['message', 'showRawPackets'],
    emits: ['callsignSelected', 'frequencySelected'],
    components: {
        ChatRxPacket,
        ChatRxMessage,
        ChatTxFrame,
    },
    methods: {
    },
    template: `
        <ChatRxPacket v-if="showRawPackets && message.Type === 'RX.ACTIVITY'" :message=message @frequencySelected="e => $emit('frequencySelected', e)" />
        <ChatRxMessage @callsignSelected="e => $emit('callsignSelected', e)" @frequencySelected="e => $emit('frequencySelected', e)" v-if="message.Type === 'RX.DIRECTED' || message.Type === 'RX.DIRECTED.ME'" :message=message :isDirectedToMe="message.Type === 'RX.DIRECTED.ME'" />
        <ChatTxFrame v-if="message.Type === 'TX.FRAME'" :message=message />
    `
}