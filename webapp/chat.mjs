import axios from 'axios'
import ChatMessage from './chat-message.mjs'

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
        }
    },
    methods: {
        chatScroll(evt) {
            const el = evt.target
            const pos = el.scrollTop / el.scrollHeight
            if (pos < 0.2 && !this.atTop && !this.loadingBefore) {
                this.loadingBefore = true
                this.fetchMessagesBefore()
            }
        },
        fetchMessages(from) {
            return axios.get('/api/rx-packets', {
                params: {
                    from: from,
                    direction: "before"
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
            console.log("TEST")
            const from = this.messages[0].Timestamp
            this.fetchMessages(from).then(result => {
                this.messages = this.messages.unshift(...result)
                this.loadingBefore = false
                this.atBottom = false
                if(!result.length) {
                    this.atTop = true
                }
            })
        },
    },
    template: `
    <div class="chat">
        <div class="chat-header clearfix">
            <div class="row">
                <div class="col-lg-6">
                    <a href="javascript:void(0);" data-toggle="modal" data-target="#view_info">
                    <img src="./avatar2.png" alt="avatar">
                    </a>
                    <div class="chat-about">
                    <h6 class="m-b-0">Aiden Chavez</h6>
                    <small>Last seen: 2 hours ago</small>
                    </div>
                </div>
            </div>
        </div>
        <div class="loader" v-if="loadingBefore">LOADING</div>
        <div class="chat-history" @scroll=chatScroll ref="chatHistory">
            <ul class="m-b-0">
                <ChatMessage v-for="message in messages" :key=message.Id :message=message />
            </ul>
        </div>
        <div class="loader" v-if="loadingAfter">LOADING</div>
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
