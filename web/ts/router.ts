import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

import Req from '../vue/req/index.vue'

Vue.use(VueRouter)

const routes = <RouteConfig[]>[
    { path: '/admin/likePostman', name: 'req', component: Req }
]

const Router = new VueRouter({
    mode: 'history',
    routes: routes
})

export default Router