import { createApp } from 'vue'
import App from './App.vue'
import router from './router/router'
import { createI18n } from 'vue-i18n'
import './main.css'


const i18n = createI18n({})

const app = createApp(App)

app.use(router)
app.use(i18n)

app.mount('#app')
