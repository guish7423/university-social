<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app"
import PcSidebar from "@/components/PcSidebar.vue"

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

<PcSidebar />

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
@media (min-width: 1024px) {
  .uni-tabbar { display: none !important; }
  .uni-tabbar-bottom { display: none !important; }
  uni-tabbar { display: none !important; }
  .uni-page-head { display: none !important; }
  .uni-page-body {
    margin-left: 220px;
    max-width: calc(100vw - 220px);
    padding-top: 0;
  }
  ::-webkit-scrollbar {
    display: block;
    width: 6px;
  }
  ::-webkit-scrollbar-thumb {
    background: #c0c4cc;
    border-radius: 3px;
  }
  ::-webkit-scrollbar-track {
    background: transparent;
  }
}
@media (min-width: 1440px) {
  .uni-page-body {
    margin-left: 220px;
    max-width: 1200px;
  }
}
</style>
