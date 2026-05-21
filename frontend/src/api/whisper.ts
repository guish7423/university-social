import { request } from "./request"

export interface WhisperData {
  id: number
  user_id: number
  content: string
  images?: string[]
  whisper_type: number
  is_anonymous: boolean
  codename?: string
  target_user_id?: number
  like_count: number
  comment_count: number
  created_at: string
  is_liked?: boolean
  author?: { nickname: string; avatar: string }
}

export interface WhisperComment {
  id: number
  whisper_id: number
  user_id: number
  content: string
  created_at: string
  codename?: string
}

export function createWhisper(data: {
  content: string
  images?: string[]
  whisper_type?: number
  target_user_id?: number
}) {
  return request<{ id: number; codename: string }>({ url: "/whispers", method: "POST", data })
}

export function listWhispers(params?: { offset?: number; limit?: number }) {
  return request<WhisperData[]>({ url: "/whispers", method: "GET", data: params as any })
}

export function getWhisper(id: number) {
  return request<WhisperData>({ url: `/whispers/${id}`, method: "GET" })
}

export function deleteWhisper(id: number) {
  return request<{ message: string }>({ url: `/whispers/${id}`, method: "DELETE" })
}

export function toggleWhisperLike(id: number) {
  return request<{ liked: boolean }>({ url: `/whispers/${id}/like`, method: "POST" })
}

export function createWhisperComment(whisperId: number, content: string) {
  return request<{ id: number }>({ url: `/whispers/${whisperId}/comments`, method: "POST", data: { content } })
}

export function listWhisperComments(whisperId: number, offset = 0, limit = 20) {
  return request<WhisperComment[]>({ url: `/whispers/${whisperId}/comments`, method: "GET", data: { offset, limit } })
}

export function whisperSummary() {
  return request<{ total: number; today: number }>({ url: "/whispers/summary", method: "GET" })
}
