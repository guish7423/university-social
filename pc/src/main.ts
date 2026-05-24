import { createApp } from "vue"
import { createPinia } from "pinia"
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import { icons } from "./icons"

import App from "./App.vue"
import router from "./router"
import "./styles/global.scss"
import "./styles/themes.css"

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

for (const [key, component] of Object.entries(icons)) {
  app.component(key, component)
}

// Global error handler
app.config.errorHandler = (err, _instance, _info) => {
  console.error("[App Error]", err)
  // Runtime errors are logged; Vue keeps rendering if possible
}

// Global promise unhandled rejection handler
window.addEventListener("unhandledrejection", (e) => {
  console.error("[Unhandled Rejection]", e.reason)
})

app.mount("#app")
