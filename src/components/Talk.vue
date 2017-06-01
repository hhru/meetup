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
            <button class="button"
                    type="button"
                    :disabled="!talk.canVote"
                    :class="{button_voted: talk.hasVoted}"
                    @click="vote(talk)">
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
            vote(talk) {
                const action = talk.hasVoted ? 'dislike' : 'like';
                this.$store.dispatch(action, talk);
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
        padding: 20px 30px;
    }

        .talk:first-child {
            margin-top: 0;
        }

    .talk__date {
        color: #8e8e8e;
        font-size: 16px;
    }

    .talk__content {
        display: inline-block;
        width: calc(100% - 68px);
        color: #666A73;
        font-size: 16px;
        line-height: 1.6;
        margin: 0;
    }

    .talk__description {
        white-space: pre-wrap;
        word-wrap: break-word;
        margin: -0.5ex 0 0;
    }

    .talk__avatar {
        display: inline-block;
        width: 48px;
        margin-right: 20px;
        vertical-align: top;
    }

    .talk__time {
        float: right;
        color: #aaa;
        font-size: 12px;
    }

    .talk__title {
        display: block;
        margin: 10px 0 30px;
        color: #282C35;
        font-size: 24px;
        font-weight: 600;
        line-height: 1.2;
        padding-bottom: 20px;
        border-bottom: 1px solid #eee;
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

        .button:hover {
            cursor: pointer;
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

        .button:disabled:after {
            content: none;
        }

        .button_voted {
            color: #7ed321;
        }

        .button_voted:after {
            background-position: 100%;
        }
</style>
