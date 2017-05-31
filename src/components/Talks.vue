<template>
    <ul class="talks">
        <talk v-for="item in talks" :key="item.key" :talk="item"/>
    </ul>
</template>

<script>
    import Talk from './Talk';

    export default {
        props: ['talks'],
        components: {
            Talk,
        },
    };
</script>

<style>
    .talks {
        margin: 0 auto;
        list-style: none;
        padding: 10px;
        max-width: 600px;
    }
</style>
