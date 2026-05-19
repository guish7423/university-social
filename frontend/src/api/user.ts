import { request } from "./request"

export interface UserInfo {
  id: number
  open_id: string
  nickname: string
  avatar: string
  phone: string
  school: string
  is_verified: boolean
  role: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export function wxLogin(code: string, nickname?: string, avatar?: string) {
  return request<LoginResponse>({
    url: "/auth/login",
    method: "POST",
    data: { code, nickname, avatar },
  })
}

export function getProfile() {
  return request<UserInfo>({
    url: "/user/profile",
    method: "GET",
  })
}

export function updateProfile(data: { nickname?: string; avatar?: string; school?: string }) {
  return request<{ message: string }>({
    url: "/user/profile",
    method: "PUT",
    data,
  })
}
