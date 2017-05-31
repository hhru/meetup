<template>
    <div>
        <div class="header-background">
            <div class="header">
                <div class="logo"></div>
            </div>
            <div class="menu">
                <router-link :to="{name: 'scheduled'}" class="menu-link">Расписание</router-link>
                <router-link :to="{name: 'proposed'}" class="menu-link">Предложенные темы</router-link>
                <router-link :to="{name: 'done'}" class="menu-link">Архив</router-link>
            </div>
        </div>
        <router-view></router-view>
    </div>
</template>

<script>
export default {
    name: 'app',
    created() {
        this.$store.dispatch('getAllTalks');
    },
};
</script>

<style>
    body {
        font-family: Helvetica, Arial, sans-serif;
        background: #f8f8fa;
        color: #666a73;
        margin: 0;
        padding: 0;
        min-width: 320px;
    }

    a {
        color: #09f;
        text-decoration: none;
        outline: none;
    }

    a:hover,
    a:active,
    a:focus {
        text-decoration: underline;
    }

    .header-background {
        position: relative;
        height: 250px;
        background: linear-gradient(
                rgba(0, 0, 0, 0.5),
                rgba(0, 0, 0, 0.5)
        ), url('./assets/background.jpg') no-repeat center;
        background-size: cover;
    }

    .header {
        max-width: 1280px;
        margin: 0 auto;
        padding: 30px;
    }

    .logo {
        float: left;
        width: 120px;
        height: 35px;
        background: url('./assets/logo.svg') no-repeat;
        background-size: 100%;
    }

    .menu {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        margin: auto;
        text-align: center;
    }

    .menu-link {
        display: inline-block;
        padding: 10px;
        color: rgba(255, 255, 255, 0.75);
        text-decoration: none;
        margin: 10px;
        outline: none;
        border: 1px solid transparent;
    }

        .menu-link:hover,
        .menu-link:active,
        .menu-link:focus,
        .menu-link_active {
            color: #fff;
            text-decoration: none;
        }

        .menu-link_active {
            border: 1px solid;
        }


    @media (max-width: 600px) {
        .menu-link {
            display: block;
        }
    }
</style>
