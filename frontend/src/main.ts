import { createApp } from 'vue'
import './assets/main.css'
import App from './App.vue'
import router, { setupRouterGuards } from './router/router'
import { VueQueryPlugin } from '@tanstack/vue-query'

import { Tabbar, TabbarItem } from 'vant'

const app = createApp(App)

app.use(Tabbar)
app.use(TabbarItem)


app.use(router)
setupRouterGuards(router)

app.use(VueQueryPlugin)

app.mount('#app')