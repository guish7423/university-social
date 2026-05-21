import { request } from "./request"

export interface DashboardStats {
  user_count: number
  post_count: number
  circle_count: number
}

export interface AdminUser {
  id: number
  nickname: string
  avatar: string
  school: string
  is_verified: boolean
  created_at: string
}

export interface AdminCircle {
  id: number
  name: string
  description: string
  member_count: number
  post_count: number
  created_at: string
}

export interface AdminReport {
  id: number
  target_type: string
  target_id: number
  reason: string
  status: string
  created_at: string
  reporter_id: number
  reporter_nickname: string
}

export function getDashboard() {
  return request<DashboardStats>({ url: "/admin/dashboard" })
}

export function listUsers(offset = 0, limit = 50) {
  return request<AdminUser[]>({ url: `/admin/users?offset=${offset}&limit=${limit}` })
}

export function banUser(id: number) {
  return request<{ message: string }>({ url: `/admin/users/${id}/ban`, method: "PUT" })
}

export function unbanUser(id: number) {
  return request<{ message: string }>({ url: `/admin/users/${id}/unban`, method: "PUT" })
}

export function listPosts(offset = 0, limit = 50) {
  return request<any[]>({ url: `/admin/posts?offset=${offset}&limit=${limit}` })
}

export function deletePost(id: number) {
  return request<{ message: string }>({ url: `/admin/posts/${id}`, method: "DELETE" })
}

export function listCircles(offset = 0, limit = 50) {
  return request<AdminCircle[]>({ url: `/admin/circles?offset=${offset}&limit=${limit}` })
}

export function deleteCircle(id: number) {
  return request<{ message: string }>({ url: `/admin/circles/${id}`, method: "DELETE" })
}

export function listReports(status = "", offset = 0, limit = 20) {
  const params = `offset=${offset}&limit=${limit}${status ? `&status=${status}` : ""}`
  return request<AdminReport[]>({ url: `/admin/reports?${params}` })
}

export function resolveReport(id: number) {
  return request<{ message: string }>({ url: `/admin/reports/${id}/resolve`, method: "PUT" })
}

export function dismissReport(id: number) {
  return request<{ message: string }>({ url: `/admin/reports/${id}/dismiss`, method: "PUT" })
}

export function listSensitive() {
  return request<string[]>({ url: "/admin/sensitive-words" })
}

export function addSensitive(word: string) {
  return request<{ message: string }>({ url: "/admin/sensitive-words", method: "POST", data: { word } })
}

export function removeSensitive(word: string) {
  return request<{ message: string }>({ url: `/admin/sensitive-words/${encodeURIComponent(word)}`, method: "DELETE" })
}
