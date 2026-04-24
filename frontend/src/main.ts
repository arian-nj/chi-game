import { createApp } from 'vue'
import './assets/main.css'
import App from './App.vue'
import { IsReleaseMode } from './lib/ReleaseMode'
import router, { setupRouterGuards } from './router/router'
import { VueQueryPlugin } from '@tanstack/vue-query'

// if (IsReleaseMode) {
//   init()
// }

const app = createApp(App)

app.use(router)
setupRouterGuards(router)

app.use(VueQueryPlugin)

app.mount('#app')