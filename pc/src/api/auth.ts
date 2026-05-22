import request from "./request"

export interface UserInfo {
  id: number
  open_id: string
  nickname: string
  avatar: string
  phone: string
  school: string
  is_verified: boolean
  role: string
  created_at: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export function devLogin(nickname: string) {
  return request.post<any, LoginResponse>("/auth/dev-login", { nickname })
}

export function getProfile() {
  return request.get<any, UserInfo>("/user/profile")
}

export function updateProfile(data: { nickname?: string; avatar?: string; school?: string }) {
  return request.put<any, { message: string }>("/user/profile", data)
}

export function getUserInfo(id: number) {
  return request.get<any, UserInfo>(`/users/${id}`)
}
