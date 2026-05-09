import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useTheme } from './composables/useTheme'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

const { initTheme } = useTheme()
initTheme()

app.mount('#app')
