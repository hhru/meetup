import Vue from 'vue';
import Router from 'vue-router';
import Scheduled from '@/components/Scheduled';
import Proposed from '@/components/Proposed';
import Done from '@/components/Done';

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            name: 'scheduled',
            component: Scheduled,
        },
        {
            path: '/proposed',
            name: 'proposed',
            component: Proposed,
        },
        {
            path: '/done',
            name: 'done',
            component: Done,
        },
    ],
    linkActiveClass: '',
    linkExactActiveClass: 'menu-link_active',
});
