<template>
    <talks :talks="talks"/>
</template>

<script>
    import {mapGetters} from 'vuex';
    import Talks from './Talks';

    export default {
        computed: mapGetters({
            talks: 'scheduled',
        }),
        components: {
            Talks,
        },
    };
</script>
