import { createApp } from "vue"
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import * as ElementPlusIconsVue from "@element-plus/icons-vue"
import DataV from "@kjgl77/datav-vue3"
import "@kjgl77/datav-vue3/dist/style.css"
import App from "./App.vue"
import "./themes.css"
import router from "./router"

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.use(ElementPlus)
app.use(DataV)
app.use(router)

app.config.errorHandler = (err, _instance, _info) => {
  console.error("[Admin Error]", err)
}
window.addEventListener("unhandledrejection", (e) => {
  console.error("[Admin Unhandled]", e.reason)
})
app.mount("#app")
