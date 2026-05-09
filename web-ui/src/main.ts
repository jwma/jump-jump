import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useTheme } from './composables/useTheme'
import i18n from './i18n'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

const { initTheme } = useTheme()
initTheme()

app.use(i18n)
app.mount('#app')
