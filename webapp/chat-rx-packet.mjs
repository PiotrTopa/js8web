export default {
    props: ['message'],
    components: {
    },
    methods: {
    },
    template: `
        <li class="clearfix packet">
            <div class="header">
                <span class="time">{{ new Date(message.Timestamp).toLocaleTimeString() }}</span>
                <span class="speed" :class="message.Speed">{{ message.Speed[0].toUpperCase() }}</span>
                <span class="freq">{{ message.Offset }}Hz</span>
            </div>
            <br />
            <div class="content">{{ message.Text }}</div>
        </li>
    `
}

                