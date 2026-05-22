<template>
  <div>
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#e6f7ff;color:#1890ff"><User /></div>
            <div>
              <div class="stat-value">{{ summary.user_count }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#fff7e6;color:#fa8c16"><Document /></div>
            <div>
              <div class="stat-value">{{ summary.post_count }}</div>
              <div class="stat-label">帖子总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#f6ffed;color:#52c41a"><Collection /></div>
            <div>
              <div class="stat-value">{{ summary.circle_count }}</div>
              <div class="stat-label">圈子总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#fff0f6;color:#eb2f96"><ChatDotRound /></div>
            <div>
              <div class="stat-value">{{ weeklyActive || '-' }}</div>
              <div class="stat-label">本周活跃</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="chart-card">
      <template #header>
        <div class="chart-header">
          <span>📊 30日趋势 (3D柱状图)</span>
          <el-radio-group v-model="chartMetric" size="small">
            <el-radio-button value="dau">DAU</el-radio-button>
            <el-radio-button value="new_users">新用户</el-radio-button>
            <el-radio-button value="new_posts">新帖子</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div ref="chartRef" style="height:400px" />
    </el-card>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header><span>📈 增量对比 (3D折线)</span></template>
          <div ref="lineChartRef" style="height:320px" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header><span>🍩 DAU 占比</span></template>
          <div ref="pieChartRef" style="height:320px" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from "vue"
import api from "@/api"
import * as echarts from "echarts"

const chartRef = ref<HTMLElement>()
const lineChartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()
const chartMetric = ref("dau")
const dailyData = ref<any[]>([])
const summary = ref({ user_count: 0, post_count: 0, circle_count: 0 })
const weeklyActive = ref(0)

const dates = computed(() => dailyData.value.map(d => d.date).reverse())
const metricData = computed(() => dailyData.value.map(d => d[chartMetric.value] || 0).reverse())

let barChart: echarts.ECharts | null = null
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

onMounted(async () => {
  try {
    const [sumRes, statsRes] = await Promise.all([
      api.get("/admin/dashboard"),
      api.get("/admin/stats/daily?days=30"),
    ])
    summary.value = sumRes.data
    dailyData.value = statsRes.data || []
    weeklyActive.value = statsRes.data?.slice(0, 7).reduce((a: number, d: any) => a + (d.dau || 0), 0) || 0
  } catch {}
  await nextTick()
  initCharts()
  resizeHandler = () => { barChart?.resize(); lineChart?.resize(); pieChart?.resize() }
  window.addEventListener("resize", resizeHandler)
})

onUnmounted(() => {
  barChart?.dispose()
  lineChart?.dispose()
  pieChart?.dispose()
  window.removeEventListener("resize", resizeHandler)
})

let resizeHandler: () => void = () => {}

watch(chartMetric, () => {
  updateBarChart()
})

function initCharts() {
  initBarChart()
  initLineChart()
  initPieChart()
}

function initBarChart() {
  if (!chartRef.value) return
  barChart = echarts.init(chartRef.value)
  updateBarChart()
}

function updateBarChart() {
  if (!barChart) return
  const nameMap: Record<string, string> = { dau: "DAU", new_users: "新用户", new_posts: "新帖子" }
  barChart.setOption({
    tooltip: { trigger: "axis" },
    grid: { left: "3%", right: "4%", bottom: "3%", containLabel: true },
    xAxis: { type: "category", data: dates.value, axisLabel: { rotate: 45, fontSize: 10 } },
    yAxis: { type: "value" },
    series: [{
      type: "bar",
      data: metricData.value.map((v: number, i: number) => ({
        value: v,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "#409eff" },
            { offset: 1, color: "#79bbff" },
          ]),
          borderRadius: [4, 4, 0, 0],
        },
      })),
      barMaxWidth: 32,
      animationDuration: 800,
      animationEasing: "cubicOut",
    }],
  })
}

function initLineChart() {
  if (!lineChartRef.value) return
  lineChart = echarts.init(lineChartRef.value)
  lineChart.setOption({
    tooltip: { trigger: "axis" },
    legend: { data: ["新用户", "新帖子", "新评论"] },
    grid: { left: "3%", right: "4%", bottom: "3%", containLabel: true },
    xAxis: { type: "category", data: dates.value, axisLabel: { rotate: 45, fontSize: 10 } },
    yAxis: { type: "value" },
    series: [
      {
        name: "新用户",
        type: "line",
        smooth: true,
        data: dailyData.value.map(d => d.new_users).reverse(),
        itemStyle: { color: "#409eff" },
        areaStyle: { color: "rgba(64,158,255,0.1)" },
      },
      {
        name: "新帖子",
        type: "line",
        smooth: true,
        data: dailyData.value.map(d => d.new_posts).reverse(),
        itemStyle: { color: "#67c23a" },
        areaStyle: { color: "rgba(103,194,58,0.1)" },
      },
      {
        name: "新评论",
        type: "line",
        smooth: true,
        data: dailyData.value.map(d => d.new_comments).reverse(),
        itemStyle: { color: "#e6a23c" },
        areaStyle: { color: "rgba(230,162,60,0.1)" },
      },
    ],
  })
}

function initPieChart() {
  if (!pieChartRef.value) return
  pieChart = echarts.init(pieChartRef.value)
  const total = dailyData.value.reduce((a: number, d: any) => a + d.new_users + d.new_posts + d.new_comments, 0) || 1
  pieChart.setOption({
    tooltip: { trigger: "item", formatter: "{b}: {c} ({d}%)" },
    series: [{
      type: "pie",
      radius: ["40%", "70%"],
      center: ["50%", "50%"],
      data: [
        { value: dailyData.value.reduce((a: number, d: any) => a + d.new_users, 0), name: "新用户" },
        { value: dailyData.value.reduce((a: number, d: any) => a + d.new_posts, 0), name: "新帖子" },
        { value: dailyData.value.reduce((a: number, d: any) => a + d.new_comments, 0), name: "新评论" },
      ],
      label: { show: true, formatter: "{b}\n{d}%" },
      itemStyle: { borderRadius: 6 },
      animationType: "scale",
    }],
  })
}
</script>

<style scoped>
.stat-row { margin-bottom: 20px; }
.stat-card { display: flex; align-items: center; gap: 16px; }
.stat-icon {
  width: 48px; height: 48px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center; font-size: 24px;
}
.stat-value { font-size: 28px; font-weight: 700; color: #303133; }
.stat-label { font-size: 13px; color: #909399; margin-top: 4px; }
.chart-card { margin-bottom: 20px; }
.chart-header { display: flex; justify-content: space-between; align-items: center; }
</style>
