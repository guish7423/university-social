import { createRouter, createWebHistory } from "vue-router"
import Dashboard from "@/views/Dashboard.vue"
import Users from "@/views/Users.vue"
import Posts from "@/views/Posts.vue"
import Circles from "@/views/Circles.vue"
import Banners from "@/views/Banners.vue"
import Announcements from "@/views/Announcements.vue"
import SensitiveWords from "@/views/SensitiveWords.vue"
import Reports from "@/views/Reports.vue"
import DataScreen from "@/views/DataScreen.vue"


const routes = [
  { path: "/", redirect: "/dashboard" },
  { path: "/dashboard", component: Dashboard, meta: { title: "数据看板" } },
  { path: "/users", component: Users, meta: { title: "用户管理" } },
  { path: "/posts", component: Posts, meta: { title: "内容管理" } },
  { path: "/circles", component: Circles, meta: { title: "圈子管理" } },
  { path: "/banners", component: Banners, meta: { title: "Banner管理" } },
  { path: "/announcements", component: Announcements, meta: { title: "公告管理" } },
  { path: "/sensitive-words", component: SensitiveWords, meta: { title: "敏感词管理" } },
  { path: "/reports", component: Reports, meta: { title: "举报管理" } },
  { path: "/datascreen", component: DataScreen, meta: { title: "数据大屏" } },
  { path: "/login", component: () => import("@/views/Login.vue"), meta: { title: "登录", hideInMenu: true } },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
