import App from '@/App.vue'
import router from '@/router.ts'
import { createApp } from 'vue'
import './assets/index.css'

import VueCookies from 'vue-cookies'

const app = createApp(App);

app.use(router)
    .use(VueCookies)
    .mount('#app')
