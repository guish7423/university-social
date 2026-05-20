import { defineStore } from "pinia"
import { ref, computed } from "vue"
import type { UserInfo } from "@/api/user"

export const useUserStore = defineStore("user", () => {
  const token = ref(uni.getStorageSync("token") || "")
	const id = ref(0)
  const nickname = ref(uni.getStorageSync("nickname") || "")
  const avatar = ref(uni.getStorageSync("avatar") || "")
  const school = ref(uni.getStorageSync("school") || "")
  const isVerified = ref(!!uni.getStorageSync("isVerified"))

  const isLogin = computed(() => !!token.value)

  function setUserInfo(info: UserInfo) {
    nickname.value = info.nickname
    id.value = info.id
    avatar.value = info.avatar
    school.value = info.school
    isVerified.value = info.is_verified
    persist()
  }

  function setToken(t: string) {
    token.value = t
    uni.setStorageSync("token", t)
  }

  function persist() {
    uni.setStorageSync("nickname", nickname.value)
    uni.setStorageSync("avatar", avatar.value)
    uni.setStorageSync("userId", id.value.toString())
    uni.setStorageSync("school", school.value)
    uni.setStorageSync("isVerified", isVerified.value ? "1" : "")
  }

  function logout() {
    id.value = 0
    token.value = ""
    nickname.value = ""
    avatar.value = ""
    school.value = ""
    isVerified.value = false
    uni.clearStorageSync()
  }

  return { id, token, nickname, avatar, school, isVerified, isLogin, setUserInfo, setToken, logout }
})
