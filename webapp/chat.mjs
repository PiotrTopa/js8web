import axios from 'axios'
import ChatMessage from './chat-message.mjs'
const EXPECTED_MESSAGES_COUNT = 100

export default {
    props: ['filter', 'showRawPackets'],
    emits: ['callsignSelected', 'frequencySelected'],
    components: {
        ChatMessage,
    },
    created() {
        this.fetchNewestMessages().then(messages => {   
            this.messages = messages
            this.atBottom = true;
            this.$nextTick(_ => {
                this.scrollToBottom()
            })
        })

        this.$nextTick(_ => window.addEventListener('event', this.event))
    },
    unmounted() {
        window.removeEventListener('event', this.event)
    },
    data() {
        return {
            messages: [],
            atTop: false,
            atBottom: false,
            loadingBefore: false,
            loadingAfter: false
        }
    },
    methods: {
        chatScroll(evt) {
            const el = evt.target
            const pos = el.scrollTop / el.scrollHeight
            if (pos < 0.2 && !this.atTop && !this.loadingBefore) {
                this.loadingBefore = true
                this.fetchMessagesBefore()
            } else if (pos > 0.8 && !this.atBottom && !this.loadingAfter) {
                this.loadingAfter = true
                this.fetchMessagesAfter()
            }
        },
        fetchMessages(startTime, direction = 'before') {
            return axios.get('/api/rx-packets', {
                params: {
                    startTime: startTime,
                    direction: direction,
                    filter: this.filter,
                }
            }).then(response => response.data)
        },
        fetchNewestMessages() {
            return this.fetchMessages(new Date(Date.now()).toISOString())
        },
        fetchMessagesBefore() {
            if (this.messages.length < 1) {
                return
            }
            const from = this.messages[0].Timestamp
            this.fetchMessages(from, 'before').then(result => {
                const existingIds = this.messages.map(e => e.Id)
                const filteredResult = result.filter(e => !existingIds.includes(e.Id))
                this.messages = filteredResult.concat(this.messages.slice(0, 2 * EXPECTED_MESSAGES_COUNT))
                this.loadingBefore = false
                this.atBottom = false
                if (result.length < EXPECTED_MESSAGES_COUNT) {
                    this.atTop = true
                }
            })
        },
        fetchMessagesAfter() {
            if (this.messages.length < 1) {
                return
            }
            const from = this.messages[this.messages.length - 1].Timestamp
            this.fetchMessages(from, 'after').then(result => {
                const existingIds = this.messages.map(e => e.Id)
                const filteredResult = result.filter(e => !existingIds.includes(e.Id))
                this.messages = this.messages.slice(-2 * EXPECTED_MESSAGES_COUNT).concat(filteredResult)
                this.loadingAfter = false
                this.atTop = false
                if (result.length < EXPECTED_MESSAGES_COUNT) {
                    this.atBottom = true
                }
            })
        },
        scrollToBottom() {
            this.$refs.chatHistory.scrollTop = this.$refs.chatHistory.scrollHeight
        },
        filterMessage(message) {
            var ret = true;
            if (this.filter) {
                if (this.filter.Callsign) {
                    const heap = (message.From + ':' + message.To).toLocaleLowerCase()
                    const needle = this.filter.Callsign.toLowerCase()
                    ret &&= heap.includes(needle)
                }

                if (this.filter.Freq && this.filter.Freq.From && this.filter.Freq.To) {
                    ret &&= message.Freq >= this.filter.Freq.From
                    ret &&= message.Freq <= this.filter.Freq.To
                }
            }
            return ret
        },
        newMessage(message) {
            if (!this.filterMessage(message)) {
                return
            }

            if (this.atBottom) {
                this.messages.push(message)
                this.$nextTick(_ => this.scrollToBottom())
            }
        },
        event(evt) {
            const event = evt.detail;
            if (event.EventType == "object" && (event.WsType == "RX.PACKET" || event.WsType == "TX.FRAME")) {
                this.newMessage(event.Event)
            }
        }
    },
    template: `
    <div class="chat">
        <div class="chat-history" @scroll=chatScroll ref="chatHistory">
            <div class="history-top" v-if="atTop">(No more messages)</div>
            <div class="loader" v-if="loadingBefore">LOADING</div>
            <ul class="m-b-0">
                <ChatMessage v-for="message in messages" :key=message.Id :message=message :showRawPackets=showRawPackets @callsignSelected="e => $emit('callsignSelected', e)" @frequencySelected="e => $emit('frequencySelected', e)" />
            </ul>
            <div class="loader" v-if="loadingAfter">LOADING</div>
            <div class="history-bottom" v-if="atBottom"><i class="bi bi-broadcast"></i> receiving <i class="bi bi-broadcast"></i></div>
        </div>
        <div class="chat-message clearfix" style="display: none">
            <div class="input-group mb-0">
            <div class="input-group-prepend">
                <span class="input-group-text"><i class="fa fa-send"></i></span>
            </div>
            <input type="text" class="form-control" placeholder="Enter text here...">
            </div>
        </div>
    </div>`
}
