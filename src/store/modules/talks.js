import talks from '../../api/talks';
import * as types from '../mutation-types';

const state = {
    all: [],
};

const PROPOSED_STATUSES = new Set(['Approved', 'Proposal']);

const getters = {
    scheduled: state => state.all.filter(item => item.status === 'Scheduled').sort((first, second) => {
        return +new Date(first.duedate) < +new Date(second.duedate) ? -1 : 1;
    }),
    proposed: state => state.all.filter(item => PROPOSED_STATUSES.has(item.status)).sort((first, second) => {
        return first.votes > second.votes ? -1 : 1;
    }),
    done: state => state.all.filter(item => item.status === 'Done').sort((first, second) => {
        return +new Date(first.duedate) > +new Date(second.duedate) ? -1 : 1;
    }),
};

const actions = {
    getAllTalks({commit}) {
        commit(types.START_REQUEST);
        talks.getAll().then((talks) => {
            commit(types.RECEIVE_TALKS, {talks});
            commit(types.FINISH_REQUEST);
        });
    },
    like({commit}, talk) {
        commit(types.START_REQUEST);
        talks.like(talk.key).then((talk) => {
            commit(types.RECEIVE_TALK, {talk});
            commit(types.FINISH_REQUEST);
        });
    },
    dislike({commit}, talk) {
        commit(types.START_REQUEST);
        talks.dislike(talk.key).then((talk) => {
            commit(types.RECEIVE_TALK, {talk});
            commit(types.FINISH_REQUEST);
        });
    },
};

const mutations = {
    [types.RECEIVE_TALKS](state, {talks}) {
        state.all = talks;
    },
    [types.RECEIVE_TALK](state, {talk}) {
        const index = state.all.findIndex(item => item.key === talk.key);
        state.all.splice(index, 1, talk);
    },
};

export default {
    state,
    getters,
    actions,
    mutations,
};
