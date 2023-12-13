import Vue from 'vue';
import ViewUI from 'view-design';
import VueI18n from 'vue-i18n'
import Router from './router'
import apiZH from './i18n/zh-CN'

import App from '../vue/app.vue'
import 'view-design/dist/styles/iview.css';

Vue.use(VueI18n)
Vue.use(ViewUI);

const i18n = new VueI18n({
    locale: 'zh-CN',
    messages: {
        'zh-CN': <any>apiZH
    }
})

new Vue({
    el: '#app',
    i18n,
    router: Router,
    render: h => h(App),
})