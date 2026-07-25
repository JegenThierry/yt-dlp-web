import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import ui from '@nuxt/ui/vue-plugin'
import router from './router'

import './main.css'

const app = createApp(App)

app.use(createPinia())
app.use(ui)
app.use(router)

app.mount('#app')
