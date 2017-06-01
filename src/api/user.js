import Vue from 'vue';

export default {
    get(request) {
        return Vue.http.get('api/me', request)
            .then(response => Promise.resolve(response.data))
            .catch(error => Promise.reject(error));
    },
};
