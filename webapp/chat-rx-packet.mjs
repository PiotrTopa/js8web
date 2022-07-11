export default {
    props: ['message'],
    components: {
    },
    methods: {
    },
    template: `
        <li class="clearfix packet">
            <div class="header">
                <span class="time">{{ message.Timestamp }}</span>
            </div>
            <div class="content">{{ message.Text }}</div>
        </li>
    `
}

                