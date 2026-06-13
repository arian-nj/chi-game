import { createHead } from '@unhead/vue/client'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router/router'
import { i18n } from './i18n'
import './main.css'

const app = createApp(App)
const head = createHead()

app.use(head)
app.use(router)
app.use(i18n)

app.mount('#app')
