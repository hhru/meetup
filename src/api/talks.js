import Vue from 'vue';

export default {
    getAll(request) {
        return Vue.http.get('api/talks', request)
            .then(response => Promise.resolve(response.data.talks))
            .catch(error => Promise.reject(error));
    },

    like(key) {
        return Vue.http.get(`api/${key}/like`)
            .then(response => Promise.resolve(response.data))
            .catch(error => Promise.reject(error));
    },
};
