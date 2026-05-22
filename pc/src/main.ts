import { createApp } from "vue"
import { createPinia } from "pinia"
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import { icons } from "./icons"

import App from "./App.vue"
import router from "./router"
import "./styles/global.scss"

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

for (const [key, component] of Object.entries(icons)) {
  app.component(key, component)
}

app.mount("#app")
