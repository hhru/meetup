import user from '../../api/user';
import * as types from '../mutation-types';

const state = {
    user: [],
};

const getters = {
    user: state => state.user,
};

const actions = {
    fetchUser({commit}) {
        commit(types.START_REQUEST);
        user.get().then((user) => {
            commit(types.RECEIVE_USER, {user});
            commit(types.FINISH_REQUEST);
        }, () => {
            window.location = '/oauth/login';
        });
    },
};

const mutations = {
    [types.RECEIVE_USER](state, {user}) {
        state.user = user;
    },
};

export default {
    state,
    getters,
    actions,
    mutations,
};
