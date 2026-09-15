import Vue from 'vue';
import ViewUI from 'view-design';
import VueI18n from 'vue-i18n'
import Router from './router'
import apiZH from './i18n/zh-CN'
import apiEN from './i18n/en-US'

import App from '../vue/app.vue'
import 'view-design/dist/styles/iview.css';

Vue.use(VueI18n)
Vue.use(ViewUI);

// 检测语言: cookie > query > navigator > 默认中文
function detectLocale(): string {
    // 从 cookie 读取
    const match = document.cookie.match(/locale=([^;]+)/)
    if (match) {
        return match[1]
    }
    // 从 URL query 读取
    const urlParams = new URLSearchParams(window.location.search)
    const lang = urlParams.get('lang')
    if (lang) {
        return lang
    }
    // 从浏览器语言检测
    const navLang = navigator.language || (navigator as any).userLanguage
    if (navLang) {
        if (navLang.startsWith('zh')) return 'zh-CN'
        if (navLang.startsWith('en')) return 'en-US'
    }
    return 'zh-CN'
}

const locale = detectLocale()

const i18n = new VueI18n({
    locale: locale,
    messages: {
        'zh-CN': <any>apiZH,
        'en-US': <any>apiEN
    }
})

new Vue({
    el: '#app',
    i18n,
    router: Router,
    render: h => h(App),
})
