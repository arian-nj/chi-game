import { createApp } from 'vue'
import './assets/main.css'

import App from './App.vue'
import router from './router/router'
import { i18n } from './i18n'

const app = createApp(App)

app.use(i18n)
app.use(router)

app.mount('#app')
