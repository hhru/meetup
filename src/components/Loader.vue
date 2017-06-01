<template>
    <div class="loadingbar"
         v-if="visible"
         v-on:transitionend="transitionEnd"
         :style="{width: width}"></div>
</template>

<script>
    export default {
        props: ['pending'],
        data() {
            return {
                inTransition: false,
                width: '1%',
            };
        },
        computed: {
            visible() {
                if (this.width === '100%' && !this.inTransition) {
                    return false;
                }
                return true;
            },
        },
        watch: {
            pending(value) {
                if (value === 0) {
                    this.width = '100%';
                } else {
                    this.width = `${(50 + Math.random() * 30) / value}%`;
                }
                this.inTransition = true;
            },
        },
        methods: {
            transitionEnd() {
                this.inTransition = false;
            },
        },
    };
</script>

<style>
    .loadingbar {
        position: fixed;
        z-index: 2147483647;
        top: 0;
        height: 3px;
        background: #d71920;
        border-radius: 1px;
        transition: all 500ms ease-in-out;
    }
</style>
