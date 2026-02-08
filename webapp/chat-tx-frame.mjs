export default {
    props: ['message'],
    template: `
        <li class="clearfix message tx-message">
            <div class="header">
                <span class="time">{{ new Date(message.Timestamp).toLocaleString() }}</span>
                <span class="gauges">
                    <span class="gauge freq"><i class="bi bi-broadcast-pin"></i> {{ message.Offset }}Hz</span>
                    <span class="gauge speed" v-if="message.Speed"><i class="bi bi-skip-end"></i><span :class="message.Speed"> {{ message.Speed[0].toUpperCase() }}</span></span>
                </span>
                <br />
                <span class="from"><i class="bi bi-send"></i> TX</span>
                <span class="selected" v-if="message.Selected"> → {{ message.Selected }}</span>
            </div>
            <div class="content">Transmitted frame</div>
        </li>
    `
}
