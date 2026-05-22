import { defineStore } from "pinia"
import { ref, computed } from "vue"
import type { UserInfo } from "@/api/auth"

export const useUserStore = defineStore("user", () => {
  const token = ref(localStorage.getItem("token") || "")
  const user = ref<UserInfo | null>(JSON.parse(localStorage.getItem("user") || "null"))

  const isLogin = computed(() => !!token.value)
  const userId = computed(() => user.value?.id || 0)
  const nickname = computed(() => user.value?.nickname || "")
  const avatar = computed(() => user.value?.avatar || "")
  const school = computed(() => user.value?.school || "")
  const isVerified = computed(() => user.value?.is_verified || false)
  const role = computed(() => user.value?.role || "user")

  function setAuth(t: string, u: UserInfo) {
    token.value = t
    user.value = u
    localStorage.setItem("token", t)
    localStorage.setItem("user", JSON.stringify(u))
  }

  function setUser(u: UserInfo) {
    user.value = u
    localStorage.setItem("user", JSON.stringify(u))
  }

  function logout() {
    token.value = ""
    user.value = null
    localStorage.removeItem("token")
    localStorage.removeItem("user")
  }

  return { token, user, isLogin, userId, nickname, avatar, school, isVerified, role, setAuth, setUser, logout }
})
