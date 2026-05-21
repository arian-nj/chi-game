import { createApp } from 'vue'
import './assets/main.css'
import App from './App.vue'
import router, { setupRouterGuards } from './router/router'
import { VueQueryPlugin } from '@tanstack/vue-query'


const app = createApp(App)


app.use(router)
setupRouterGuards(router)

app.use(VueQueryPlugin)

app.mount('#app')
