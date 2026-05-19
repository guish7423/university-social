import { useUserStore } from "@/store/user"

const BASE_URL = import.meta.env.VITE_API_BASE || "http://localhost:8080/api/v1"

interface RequestOptions {
  url: string
  method?: "GET" | "POST" | "PUT" | "DELETE"
  data?: Record<string, any>
  header?: Record<string, string>
}

export function request<T = any>(options: RequestOptions): Promise<T> {
  const userStore = useUserStore()
  const header: Record<string, string> = {
    "Content-Type": "application/json",
    ...options.header,
  }
  if (userStore.token) {
    header["Authorization"] = `Bearer ${userStore.token}`
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: `${BASE_URL}${options.url}`,
      method: options.method || "GET",
      data: options.data,
      header,
      success: (res) => {
        if (res.statusCode === 401) {
          userStore.logout()
          uni.reLaunch({ url: "/pages/login/index" })
          reject(new Error("未登录"))
          return
        }
        resolve(res.data as T)
      },
      fail: (err) => {
        uni.showToast({ title: "网络错误", icon: "none" })
        reject(err)
      },
    })
  })
}
