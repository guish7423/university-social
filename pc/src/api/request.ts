import axios from "axios"
import { ElMessage } from "element-plus"
import { useUserStore } from "@/stores/user"
import router from "@/router"

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem("token")
  if (token) {
    config.headers.Authorization = "Bearer " + token
  }
  return config
})

request.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      const store = useUserStore()
      store.logout()
      router.push("/login")
    }
    const msg = err.response?.data?.error || err.message || "网络错误"
    ElMessage.error(msg)
    return Promise.reject(err)
  },
)

export default request
