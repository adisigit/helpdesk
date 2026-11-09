import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import { configureHttp } from './api/http'
import { useAuthStore } from './stores/auth'
import { router } from './router'

const app = createApp(App)
app.use(createPinia())

configureHttp({
  onUnauthorized: () => {
    useAuthStore().clear()
    router.replace({
      name: 'login',
      query: { redirect: router.currentRoute.value.fullPath, expired: '1' },
    })
  },
})

app.use(router)
app.mount('#app')
