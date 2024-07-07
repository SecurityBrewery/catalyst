import './assets/main.css'
import 'easymde/dist/easymde.min.css'

import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import VueEasymde from 'vue3-easymde'

const pinia = createPinia()
const app = createApp(App)

app.use(VueEasymde)
app.use(pinia)
app.use(router)
app.use(VueQueryPlugin)

app.mount('#app')
