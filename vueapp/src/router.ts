import Main from "@/components/Main.vue"

import { createRouter, createWebHistory } from 'vue-router'

export default createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: Main,
        },
    ],
})


