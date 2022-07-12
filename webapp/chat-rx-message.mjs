export default {
    props: ['message'],
    components: {
    },
    methods: {
    },
    template: `
        <li class="clearfix message my-message">
            <div class="header">
                <span class="time">{{ new Date(message.Timestamp).toLocaleString() }}</span>
                <span class="speed" :class="message.Speed">{{ message.Speed[0].toUpperCase() }}</span>
                <span class="freq">{{ message.Offset }}Hz</span>
                <br />

                <span class="from">{{ message.From }}</span>
                <span class="grid" v-if=message.Grid>{{ message.Grid }}</span>
            </div>
            <div class="content">{{ message.Text }}</div>
        </li>
    `
}

                