import Vue from 'vue';
import Vuex from 'vuex';
import talks from './modules/talks';
import user from './modules/user';
import * as types from './mutation-types';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

const state = {
    pending: 0,
};

const getters = {
    pending: state => state.pending,
};

const mutations = {
    [types.FINISH_REQUEST](state) {
        state.pending -= 1;
    },
    [types.START_REQUEST](state) {
        state.pending += 1;
    },
};

export default new Vuex.Store({
    state,
    getters,
    mutations,
    modules: {
        talks,
        user,
    },
    strict: debug,
});
