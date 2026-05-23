<template>
  <router-view v-slot="{ Component, route }">
    <Transition name="page" mode="out-in">
      <component :is="Component" :key="route.path" />
    </Transition>
  </router-view>
</template>

<script setup lang="ts">
import { onMounted } from "vue"
import { useUserStore } from "@/stores/user"
import { getProfile } from "@/api/auth"
import router from "@/router"

const userStore = useUserStore()

onMounted(async () => {
  if (!userStore.token) return
  try {
    const profile = await getProfile()
    userStore.setUser(profile)
  } catch {
    userStore.logout()
    router.push("/login")
  }
})
</script>
