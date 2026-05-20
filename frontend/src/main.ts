import { createSSRApp } from "vue"
import { createPinia } from "pinia"
import "uno.css"
import App from "./App.vue"
import uViewPro from "uview-pro"

export function createApp() {
  const app = createSSRApp(App)
  app.use(createPinia())
  app.use(uViewPro)
  return { app }
}
