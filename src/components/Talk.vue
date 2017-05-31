<template>
    <li class="talk">
        <span v-if="talk.duedate" class="talk__date">{{ [talk.duedate, "YYYY-MM-DD"] | moment("LL") }}</span>
        <span v-if="talk.duedate" class="talk__time">{{ [talk.duedate, "YYYY-MM-DD"] | moment("from", "now") }}</span>
        <a :href="url" class="talk__title" target="_blank">{{ talk.summary }}</a>
        <div>
            <img class="talk__avatar" :src="talk.author.avatar"/><!--
             --><div class="talk__content">
            <p class="talk__description">{{ talk.description }}</p>
            <ul class="attachments" v-if="talk.attachments.length">
                <li class="attachment" v-for="attachment in talk.attachments">
                    <a :href="attachment.url" target="_blank">{{ attachment.name }}</a>
                </li>
            </ul>
        </div>
        </div>
        <div class="buttons">
            <button class="button" @click="like(talk)" v-bind:class="{button_voted: talk.hasVoted}" type="button">
                {{ talk.votes }}
            </button>
        </div>
    </li>
</template>

<script>
    export default {
        props: ['talk'],
        name: 'talk',
        computed: {
            url() {
                return `https://jira.hh.ru/browse/${this.talk.key}`;
            },
        },
        methods: {
            like(talk) {
                this.$store.dispatch('like', talk);
            },
        },
    };
</script>

<style>
    .talk {
        display: block;
        list-style: none;
        background: #fff;
        box-shadow: 0 1px 0 0 rgba(0, 0, 0, .1);
        margin-top: 10px;
        padding: 10px 20px;
    }

        .talk:first-child {
            margin-top: 0;
        }

    .talk__date {
        color: #45494e;
        font-size: 12px;
        text-transform: uppercase;
        letter-spacing: 1px;
    }

    .talk__content {
        display: inline-block;
        width: calc(100% - 68px);
        color: #666A73;
        font-size: 13px;
        margin: 0;
    }

    .talk__description {
        white-space: pre-wrap;
        word-wrap: break-word;
        margin: 0;
    }

    .talk__avatar {
        display: inline-block;
        width: 48px;
        margin-right: 20px;
        vertical-align: top;
    }

    .talk__time {
        float: right;
        color: #999;
        font-size: 11px;
    }

    .talk__title {
        display: block;
        margin: 10px 0;
        color: #282C35;
        font-size: 15px;
        font-weight: 600;
    }

    .attachments {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .attachment {
        margin: 5px 0;
    }

    .buttons {
        text-align: right;
    }

    .button {
        position: relative;
        height: 32px;
        padding: 0;
        background: none;
        border: none;
        outline: none;
        vertical-align: middle;
        line-height: 32px;
        font-weight: bold;
    }

        .button:after {
            content: '';
            display: inline-block;
            width: 20px;
            height: 20px;
            vertical-align: middle;
            margin-left: 5px;
            background: url('../assets/vote.svg') no-repeat;
            background-size: 200%;
        }

        .button_voted {
            color: #7ed321;
        }

        .button_voted:after {
            background-position: 100%;
        }
</style>
