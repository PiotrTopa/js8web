
export default {
    props: ['packets'],
    methods: {
    },
    template: `
    <div v-for="packet in packets">
        {{ packet }}
    </div>
    `
}
