import { VueQueryPlugin } from '@tanstack/vue-query'
import { createApp } from 'vue'
import { createHead } from '@unhead/vue/client'
import './assets/main.css'

import App from './App.vue'
import { i18n } from './i18n'
import router from './router/router'
import { queryClient } from './libs/vue-query'

const app = createApp(App)
const head = createHead()

app.use(head)
app.use(i18n)
app.use(router)
app.use(VueQueryPlugin, { queryClient })

app.mount('#app')
