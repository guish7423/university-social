<template>
  <dv-full-screen-container>
    <div class="screen-bg">
      <!-- Header -->
      <div class="header">
        <dv-decoration-5 class="header-decoration" />
        <div class="header-content">
          <dv-decoration-8 class="header-left-dec" />
          <div class="header-title">
            <h1>校园社 数据大屏</h1>
            <p>洛阳高校 · 实时运营监控平台</p>
          </div>
          <dv-decoration-8 class="header-right-dec" :reverse="true" />
        </div>
        <dv-decoration-5 class="header-decoration" :reverse="true" />
      </div>

      <!-- Main Content -->
      <div class="main-content">
        <!-- Left Column -->
        <div class="column left-column">
          <dv-border-box-13 class="panel">
            <div class="panel-header">
              <dv-decoration-6 class="panel-decoration" />
              <span>实时活跃统计</span>
              <dv-decoration-6 class="panel-decoration" :reverse="true" />
            </div>
            <div class="panel-body">
              <div class="stat-grid">
                <div class="stat-item">
                  <dv-border-box-12>
                    <dv-digital-flop :config="onlineUsers" style="width:100%;height:60px" />
                    <div class="stat-label">在线用户</div>
                  </dv-border-box-12>
                </div>
                <div class="stat-item">
                  <dv-border-box-12>
                    <dv-digital-flop :config="dailyActive" style="width:100%;height:60px" />
                    <div class="stat-label">日活跃</div>
                  </dv-border-box-12>
                </div>
              </div>
              <div class="stat-grid" style="margin-top:12px">
                <div class="stat-item">
                  <dv-border-box-12>
                    <dv-digital-flop :config="weeklyPosts" style="width:100%;height:60px" />
                    <div class="stat-label">本周帖子</div>
                  </dv-border-box-12>
                </div>
                <div class="stat-item">
                  <dv-border-box-12>
                    <dv-digital-flop :config="totalUsers" style="width:100%;height:60px" />
                    <div class="stat-label">总用户数</div>
                  </dv-border-box-12>
                </div>
              </div>
            </div>
          </dv-border-box-13>

          <dv-border-box-13 class="panel" style="margin-top:16px">
            <div class="panel-header">
              <dv-decoration-6 class="panel-decoration" />
              <span>活跃度排行</span>
              <dv-decoration-6 class="panel-decoration" :reverse="true" />
            </div>
            <div class="panel-body">
              <dv-scroll-ranking-board :config="rankingConfig" style="width:100%;height:260px" />
            </div>
          </dv-border-box-13>
        </div>

        <!-- Center Column -->
        <div class="column center-column">
          <div class="map-container">
            <dv-border-box-8 class="map-box">
              <template #default>
                <dv-flyline-chart-enhanced :config="flylineConfig" style="width:100%;height:100%" />
              </template>
            </dv-border-box-8>
          </div>
          <div class="bottom-charts">
            <dv-border-box-13 class="chart-panel">
              <div class="chart-header">
                <dv-decoration-11>内容数据趋势</dv-decoration-11>
              </div>
              <dv-active-ring-chart :config="activeRingConfig" style="width:100%;height:200px" />
            </dv-border-box-13>
            <dv-border-box-13 class="chart-panel" style="margin-left:12px">
              <div class="chart-header">
                <dv-decoration-11>板块分布</dv-decoration-11>
              </div>
              <dv-capsule-chart :config="capsuleConfig" style="width:100%;height:200px" />
            </dv-border-box-13>
          </div>
        </div>

        <!-- Right Column -->
        <div class="column right-column">
          <dv-border-box-13 class="panel">
            <div class="panel-header">
              <dv-decoration-6 class="panel-decoration" />
              <span>实时动态</span>
              <dv-decoration-6 class="panel-decoration" :reverse="true" />
            </div>
            <div class="panel-body">
              <dv-scroll-board :config="scrollConfig" style="width:100%;height:280px" />
            </div>
          </dv-border-box-13>

          <dv-border-box-13 class="panel" style="margin-top:16px">
            <div class="panel-header">
              <dv-decoration-6 class="panel-decoration" />
              <span>数据概览</span>
              <dv-decoration-6 class="panel-decoration" :reverse="true" />
            </div>
            <div class="panel-body">
              <div class="overview-grid">
                <div class="overview-item">
                  <dv-water-level-pond :config="waterConfig" style="width:100%;height:100%" />
                  <div class="overview-label">服务器负载</div>
                </div>
                <div class="overview-item">
                  <dv-percent-pond :config="percentConfig" style="width:100%;height:40px" />
                  <div class="overview-label">存储使用率</div>
                </div>
                <div class="overview-item">
                  <dv-conical-column-chart :config="conicalConfig" style="width:100%;height:100%" />
                  <div class="overview-label">今日新增</div>
                </div>
              </div>
            </div>
          </dv-border-box-13>
        </div>
      </div>
    </div>
  </dv-full-screen-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import api from "@/api"

const onlineUsers = ref({
  number: [2856, 3120, 3456, 3289, 4012],
  content: "{nt}"
})

const dailyActive = ref({
  number: [12056, 13520, 15280, 14890, 16750],
  content: "{nt}"
})

const weeklyPosts = ref({
  number: [3520, 4280, 3960, 5130, 5720],
  content: "{nt}"
})

const totalUsers = ref({
  number: [85600, 92300, 105420, 118560, 128900],
  content: "{nt}"
})

const rankingConfig = ref({
  data: [
    { name: "洛阳理工大学", value: 5872 },
    { name: "河南科技大学", value: 4521 },
    { name: "洛阳师范学院", value: 3890 },
    { name: "洛阳理工学院", value: 3245 },
    { name: "河南推拿职业学院", value: 2768 },
    { name: "洛阳职业技术学院", value: 2134 },
    { name: "洛阳科技职业学院", value: 1876 },
    { name: "郑州大学洛阳校区", value: 1456 },
  ],
  unit: "人",
})

const flylineConfig = ref({
  center: [112.45, 34.62],
  points: [
    { position: [112.45, 34.62], text: "洛阳" },
    { position: [113.54, 34.72], text: "郑州" },
    { position: [111.0, 34.0], text: "南阳" },
    { position: [114.41, 34.36], text: "开封" },
    { position: [112.2, 33.0], text: "平顶山" },
    { position: [113.63, 34.74], text: "新乡" },
    { position: [114.5, 35.3], text: "安阳" },
    { position: [112.9, 35.4], text: "济源" },
  ],
  lines: [
    { source: "洛阳", target: "郑州" },
    { source: "洛阳", target: "南阳" },
    { source: "洛阳", target: "开封" },
    { source: "洛阳", target: "平顶山" },
    { source: "洛阳", target: "新乡" },
    { source: "洛阳", target: "安阳" },
    { source: "洛阳", target: "济源" },
  ],
  lineWidth: 2,
  lineColor: "#1a98fc",
  duration: 3000,
})

const activeRingConfig = ref({
  series: [
    { name: "帖子", value: 3520 },
    { name: "动态", value: 2180 },
    { name: "评论", value: 8920 },
    { name: "点赞", value: 15600 },
    { name: "分享", value: 3240 },
  ],
})

const capsuleConfig = ref({
  data: [
    { name: "校园广场", value: 4560 },
    { name: "失物招领", value: 2150 },
    { name: "二手交易", value: 1980 },
    { name: "课程讨论", value: 1560 },
    { name: "社团活动", value: 1230 },
    { name: "校园求助", value: 980 },
  ],
  unit: "条",
})

const scrollConfig = ref({
  header: ["时间", "用户", "操作"],
  data: [
    ["14:32:15", "张三", "发布了新帖子"],
    ["14:28:40", "李四", "评论了校园活动"],
    ["14:25:12", "王五", "点赞了二手交易"],
    ["14:20:33", "赵六", "加入了跑步社团"],
    ["14:15:48", "孙七", "提交了失物招领"],
    ["14:10:22", "周八", "发布了课程笔记"],
    ["14:05:56", "吴九", "回复了求助帖"],
    ["14:00:18", "郑十", "分享了校园资讯"],
    ["13:55:42", "刘一", "创建了新圈子"],
    ["13:50:10", "陈二", "上传了活动照片"],
  ],
  rowNum: 8,
  headerBGC: "rgba(26,152,252,0.15)",
  headerHeight: 35,
})

const waterConfig = ref({
  data: [56],
  shape: "round",
  waveNum: 3,
  colors: ["#1a98fc"],
})

const percentConfig = ref({
  value: 68,
  borderWidth: 6,
  colors: [{ color: "#1a98fc", percent: 100 }],
})

const conicalConfig = ref({
  data: [
    { name: "用户注册", value: 156 },
    { name: "帖子发布", value: 89 },
    { name: "圈子创建", value: 12 },
    { name: "交易完成", value: 45 },
  ],
})

const summary = ref({ user_count: 0, post_count: 0, circle_count: 0 })
const dailyStats = ref<any[]>([])

onMounted(async () => {
  try {
    const [sumRes, statsRes] = await Promise.all([
      api.get("/admin/dashboard"),
      api.get("/admin/stats/daily?days=30"),
    ])
    summary.value = sumRes.data
    dailyStats.value = statsRes.data || []

    const s = summary.value
    totalUsers.value.number = [s.user_count * 1.2, s.user_count]
    dailyActive.value.number = [dailyStats.value[0]?.dau || dailyActive.value.number[0]]
    onlineUsers.value.number = [Math.floor((dailyStats.value[0]?.dau || 0) * 0.3)]
    weeklyPosts.value.number = [dailyStats.value.slice(0,7).reduce((a:number,d:any)=>a+(d.new_posts||0),0)]
  } catch {}
})
</script>

<style>
.screen-bg {
  width: 100%;
  height: 100%;
  background: #020b2a;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  color: #d3d6dd;
}

.header {
  height: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  padding-top: 8px;
}

.header-decoration {
  width: 60%;
  height: 12px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 8px 0;
}

.header-title {
  text-align: center;
  margin: 0 24px;
}

.header-title h1 {
  margin: 0;
  font-size: 36px;
  font-weight: 700;
  background: linear-gradient(135deg, #1a98fc, #00d4ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 6px;
}

.header-title p {
  margin: 4px 0 0;
  font-size: 14px;
  color: rgba(255,255,255,0.5);
  letter-spacing: 3px;
}

.main-content {
  flex: 1;
  display: flex;
  padding: 12px 16px;
  gap: 16px;
  min-height: 0;
}

.column {
  display: flex;
  flex-direction: column;
}

.left-column, .right-column {
  width: 22%;
  min-width: 240px;
}

.center-column {
  flex: 1;
  min-width: 0;
}

.panel {
  flex: 1;
  min-height: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #d3d6dd;
}

.panel-decoration {
  width: 60px;
  height: 20px;
}

.panel-body {
  padding: 8px 12px;
}

.stat-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.stat-item {
  text-align: center;
}

.stat-item .dv-digital-flop {
  font-size: 28px;
  font-weight: 700;
}

.stat-label {
  font-size: 12px;
  color: rgba(255,255,255,0.6);
  margin-top: 2px;
  padding-bottom: 8px;
}

.map-container {
  flex: 1;
  min-height: 0;
}

.map-box {
  height: 100%;
  padding: 4px;
}

.map-box :deep(.border-box-content) {
  height: 100%;
}

.bottom-charts {
  display: flex;
  margin-top: 16px;
  height: 260px;
  flex-shrink: 0;
}

.chart-panel {
  flex: 1;
  min-width: 0;
  padding: 8px;
}

.chart-header {
  text-align: center;
  font-size: 14px;
  color: #d3d6dd;
  margin-bottom: 4px;
}

.overview-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.overview-item {
  text-align: center;
  min-height: 60px;
}

.overview-label {
  font-size: 12px;
  color: rgba(255,255,255,0.6);
  margin-top: 4px;
}
</style>
