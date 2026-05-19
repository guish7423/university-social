import { createRouter, createWebHistory } from "vue-router"
import Dashboard from "@/views/Dashboard.vue"
import Users from "@/views/Users.vue"
import Posts from "@/views/Posts.vue"
import Circles from "@/views/Circles.vue"

const routes = [
  { path: "/", redirect: "/dashboard" },
  { path: "/dashboard", component: Dashboard, meta: { title: "数据看板" } },
  { path: "/users", component: Users, meta: { title: "用户管理" } },
  { path: "/posts", component: Posts, meta: { title: "内容管理" } },
  { path: "/circles", component: Circles, meta: { title: "圈子管理" } },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
