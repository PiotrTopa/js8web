export default {
    props: ['message'],
    components: {
    },
    methods: {
    },
    template: `
        <li class="clearfix message my-message">
            <div class="header">
                <span class="time">{{ message.Timestamp }}</span>
                <span class="from">{{ message.From }}</span>
                <span class="to" v-if="message.To">{{ message.To }}</span>
            </div>
            <div class="content">{{ message.Text }}</div>
        </li>
    `
}

                