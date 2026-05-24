import { createRouter, createWebHashHistory } from "vue-router"
import MainLayout from "@/layouts/MainLayout.vue"
import { useUserStore } from "@/stores/user"

const router = createRouter({
  history: createWebHashHistory(),
  routes: [

    {
      path: "/login",
      name: "Login",
      component: () => import("@/pages/login/index.vue"),
    },
    {
      path: "/",
      component: MainLayout,
      redirect: "/home",
      children: [
        { path: "home", name: "Home", meta: { title: "首页" }, component: () => import("@/pages/home/index.vue") },
        { path: "square", name: "Square", meta: { title: "广场" }, component: () => import("@/pages/square/index.vue") },
        { path: "search", name: "Search", meta: { title: "搜索" }, component: () => import("@/pages/search/index.vue") },
        { path: "post/:id", name: "PostDetail", meta: { title: "动态详情" }, component: () => import("@/pages/post/detail.vue") },
        { path: "post/create", name: "CreatePost", meta: { title: "发布动态" }, component: () => import("@/pages/post/create.vue") },
        { path: "circles", name: "Circles", meta: { title: "圈子" }, component: () => import("@/pages/circle/list.vue") },
        { path: "circles/:id", name: "CircleDetail", meta: { title: "圈子详情" }, component: () => import("@/pages/circle/detail.vue") },
        { path: "circle/create", name: "CreateCircle", meta: { title: "创建圈子" }, component: () => import("@/pages/circle/create.vue") },
        { path: "activities", name: "Activities", meta: { title: "活动" }, component: () => import("@/pages/activity/list.vue") },
        { path: "activities/:id", name: "ActivityDetail", meta: { title: "活动详情" }, component: () => import("@/pages/activity/detail.vue") },
        { path: "activity/create", name: "CreateActivity", meta: { title: "创建活动" }, component: () => import("@/pages/activity/create.vue") },
        { path: "products", name: "Products", meta: { title: "二手" }, component: () => import("@/pages/product/list.vue") },
        { path: "products/:id", name: "ProductDetail", meta: { title: "商品详情" }, component: () => import("@/pages/product/detail.vue") },
        { path: "product/create", name: "CreateProduct", meta: { title: "发布商品" }, component: () => import("@/pages/product/create.vue") },
        { path: "found", name: "Found", meta: { title: "失物招领" }, component: () => import("@/pages/found/list.vue") },
        { path: "found/:id", name: "FoundDetail", meta: { title: "失物详情" }, component: () => import("@/pages/found/detail.vue") },
        { path: "found/create", name: "CreateFound", meta: { title: "发布失物" }, component: () => import("@/pages/found/create.vue") },
        { path: "friends", name: "Friends", meta: { title: "好友" }, component: () => import("@/pages/friend/index.vue") },
        { path: "chat", name: "Chat", meta: { title: "聊天" }, component: () => import("@/pages/chat/index.vue") },
        { path: "chat/:userId", name: "ChatDetail", meta: { title: "聊天" }, component: () => import("@/pages/chat/detail.vue") },
        { path: "notifications", name: "Notifications", meta: { title: "通知" }, component: () => import("@/pages/notification/index.vue") },
        { path: "user/:id", name: "UserDetail", meta: { title: "个人主页" }, component: () => import("@/pages/user/detail.vue") },
        { path: "user/:id/posts", name: "UserPosts", meta: { title: "Ta的动态" }, component: () => import("@/pages/user/posts.vue") },
        { path: "user/:id/followers", name: "UserFollowers", meta: { title: "粉丝" }, component: () => import("@/pages/user/followers.vue") },
        { path: "user/:id/following", name: "UserFollowing", meta: { title: "关注" }, component: () => import("@/pages/user/following.vue") },
        { path: "profile/edit", name: "EditProfile", meta: { title: "编辑资料" }, component: () => import("@/pages/user/edit.vue") },
        { path: "account/security", name: "AccountSecurity", meta: { title: "账号安全" }, component: () => import("@/pages/account/security.vue") },
        { path: "points", name: "Points", meta: { title: "积分" }, component: () => import("@/pages/points/index.vue") },
        { path: "favorites", name: "Favorites", meta: { title: "收藏" }, component: () => import("@/pages/favorite/index.vue") },
        { path: "whispers", name: "Whispers", meta: { title: "树洞" }, component: () => import("@/pages/whisper/index.vue") },
        { path: "whispers/:id", name: "WhisperDetail", meta: { title: "树洞详情" }, component: () => import("@/pages/whisper/detail.vue") },
        { path: "courses", name: "Courses", meta: { title: "课程" }, component: () => import("@/pages/course/index.vue") },
        { path: "courses/:id", name: "CourseDetail", meta: { title: "课程详情" }, component: () => import("@/pages/course/detail.vue") },
        { path: "verification", name: "Verification", meta: { title: "认证" }, component: () => import("@/pages/verification/index.vue") },
        { path: "campus/calendar", name: "CampusCalendar", meta: { title: "校历" }, component: () => import("@/pages/campus/calendar.vue") },
        { path: "campus/directory", name: "CampusDirectory", meta: { title: "通讯录" }, component: () => import("@/pages/campus/directory.vue") },
        { path: "campus/rooms", name: "CampusRooms", meta: { title: "空教室" }, component: () => import("@/pages/campus/rooms.vue") },
        { path: "campus", redirect: "/campus/calendar" },
        { path: "topics", name: "Topics", meta: { title: "话题" }, component: () => import("@/pages/topic/index.vue") },
        { path: "reviews", name: "Reviews", meta: { title: "评价" }, component: () => import("@/pages/review/index.vue") },
        { path: "reviews/create", name: "CreateReview", meta: { title: "发起评价" }, component: () => import("@/pages/review/create.vue") },
        { path: "reviews/:id", name: "ReviewDetail", meta: { title: "评价详情" }, component: () => import("@/pages/review/detail.vue") },

        { path: "announcements", name: "Announcements", meta: { title: "公告" }, component: () => import("@/pages/announcement/index.vue") },
        { path: "circles/:id/members", name: "CircleMembers", meta: { title: "成员" }, component: () => import("@/pages/circle/members.vue") },
        { path: "university/:id", name: "UniversityDetail", meta: { title: "高校详情" }, component: () => import("@/pages/university/detail.vue") },

        { path: "invite", name: "Invite", meta: { title: "邀请码" }, component: () => import("@/pages/invite/index.vue") },
        { path: "university", name: "Universities", meta: { title: "高校" }, component: () => import("@/pages/university/index.vue") },
        { path: "/:pathMatch(.*)*", name: "NotFound", component: () => import("@/pages/error/NotFound.vue") },
      ],
    },
  ],
  scrollBehavior(to: any, from: any, savedPosition: any) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, behavior: "smooth" as const }
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const token = localStorage.getItem("token")
  if (to.name !== "Login" && !token) {
    return { name: "Login" }
  }
  const userStore = useUserStore()
  if (userStore.userId) {
    document.title = (to.meta.title as string) ? `${to.meta.title} · 校园社` : "校园社"
  }
})

export default router
