import axios from 'axios'
import ChatMessage from './chat-message.mjs'
const EXPECTED_MESSAGES_COUNT = 100

function uidGenerator() {
    var S4 = function() {
       return (((1+Math.random())*0x10000)|0).toString(16).substring(1);
    };
    return (S4()+S4()+"-"+S4()+"-"+S4()+"-"+S4()+"-"+S4()+S4()+S4());
}

export default {
    props: ['filter'],
    components: {
        ChatMessage,
    },
    created() {
        this.fetchNewestMessages().then(messages => {
            this.messages = messages
            this.$nextTick(_ => {
                this.$refs.chatHistory.scrollTop = this.$refs.chatHistory.scrollHeight
            })
        })
    },
    data() {
        return {
            messages: [],
            atTop: false,
            atBottom: false,
            loadingBefore: false,
            loadingAfter: false,
            showRawPackets: true,
            uid: uidGenerator(),
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
        fetchMessages(from, direction = 'before') {
            return axios.get('/api/rx-packets', {
                params: {
                    from: from,
                    direction: direction
                }
            }).then(response => {
                return new Promise((resolve, reject) => {
                    resolve(response.data)
                })
            })
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
    },
    template: `
    <div class="chat">
        <div class="chat-header clearfix">
            <div class="row">
                <div class="col-lg-6">
                    <div class="chat-about">
                        <h6 class="m-b-0">All messages</h6>
                    </div>
                    <div class="form-check form-switch settings">
                        <input class="form-check-input" type="checkbox" role="switch" :id="this.uid+'-show-raw-packets'" v-model="this.showRawPackets">
                        <label class="form-check-label" :for="this.uid+'-show-raw-packets'">Show raw packets</label>
                    </div>
                    
                </div>
            </div>
        </div>
        <div class="chat-history" @scroll=chatScroll ref="chatHistory">
            <div class="history-top" v-if="atTop">(No more messages)</div>
            <div class="loader" v-if="loadingBefore">LOADING</div>
            <ul class="m-b-0">
                <ChatMessage v-for="message in messages" :key=message.Id :message=message :showRawPackets=showRawPackets />
            </ul>
            <div class="loader" v-if="loadingAfter">LOADING</div>
            <div class="history-top" v-if="atBottom">(receiving)</div>
        </div>
        <div class="chat-message clearfix">
            <div class="input-group mb-0">
            <div class="input-group-prepend">
                <span class="input-group-text"><i class="fa fa-send"></i></span>
            </div>
            <input type="text" class="form-control" placeholder="Enter text here...">
            </div>
        </div>
    </div>`
}
