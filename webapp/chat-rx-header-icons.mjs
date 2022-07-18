var SNR_COLOR_1 = {
    red: 0, green: 0, blue: 125
};
var SNR_COLOR_2 = {
    red: 255, green: 255, blue: 0
};
var SNR_COLOR_3 = {
    red: 255, green: 0, blue: 0
};

export default {
    props: ['message'],
    emits: ['frequencySelected'],
    components: {
    },
    methods: {
        snrColor(snr) {
            var snrAligned = Math.max(-30, Math.min(20, snr)) + 30
            var fade = snrAligned / 50.0;
            var color1 = SNR_COLOR_1;
            var color2 = SNR_COLOR_2;

            fade = fade * 2;
            if (fade >= 1) {
                fade -= 1;
                color1 = SNR_COLOR_2;
                color2 = SNR_COLOR_3;
            }

            var diffRed = color2.red - color1.red;
            var diffGreen = color2.green - color1.green;
            var diffBlue = color2.blue - color1.blue;

            var gradient = {
                red: parseInt(Math.floor(color1.red + (diffRed * fade)), 10),
                green: parseInt(Math.floor(color1.green + (diffGreen * fade)), 10),
                blue: parseInt(Math.floor(color1.blue + (diffBlue * fade)), 10),
            };

            return 'rgb(' + gradient.red + ',' + gradient.green + ',' + gradient.blue + ')';
        }
    },
    template: `
        <span class="gauges">
            <span class="gauge freq"><a class="btn btn-light btn-sm" href="#" @click="$emit('frequencySelected', message.Freq)"><i class="bi bi-broadcast-pin"></i> {{ message.Offset }}Hz</a></span>
            <span class="gauge snr"><i class="bi bi-speedometer2"></i><span :style="'color: ' + snrColor(message.Snr)">{{ message.Snr > 0 ? '+' : '' }} {{ message.Snr }}</span></span>
            <span class="gauge speed" v-if="message.Speed"><i class="bi bi-skip-end"></i><span :class="message.Speed"> {{ message.Speed[0].toUpperCase() }}</span></span>
            <span class="gauge timedritft"><i class="bi bi-stopwatch"></i> {{ message.TimeDrift > 0 ? '+' : '' }}{{ message.TimeDrift }}ms</span>
        </span>
    `
}

