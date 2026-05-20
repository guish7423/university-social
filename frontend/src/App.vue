<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app"

onLaunch(() => {
  console.log("App Launch")
})
onShow(() => {
  console.log("App Show")
})
onHide(() => {
  console.log("App Hide")
})
</script>

<template>
  <view />
</template>

<style>
page {
  background-color: #f0f2f5;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif;
  overscroll-behavior: none;
  -webkit-font-smoothing: antialiased;
}

view {
  box-sizing: border-box;
}

::-webkit-scrollbar {
  display: none;
}
</style>
