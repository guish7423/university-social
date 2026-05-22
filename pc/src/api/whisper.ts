import request from "./request"

export interface WhisperData {
  id: number
  user_id: number
  content: string
  like_count: number
  comment_count: number
  is_anonymous: boolean
  created_at: string
  author?: { nickname: string; avatar: string }
  is_liked?: boolean
}

export function listWhispers(offset = 0, limit = 20) {
  return request.get<any, WhisperData[]>("/whispers", { params: { offset, limit } })
}

export function getWhisper(id: number) {
  return request.get<any, WhisperData>(`/whispers/${id}`)
}

export function createWhisper(data: { content: string; is_anonymous?: boolean }) {
  return request.post<any, { id: number }>("/whispers", data)
}

export function deleteWhisper(id: number) {
  return request.delete<any, { message: string }>(`/whispers/${id}`)
}

export function toggleWhisperLike(id: number) {
  return request.post<any, { liked: boolean }>(`/whispers/${id}/like`)
}

export function createWhisperComment(whisperId: number, content: string) {
  return request.post<any, { id: number }>(`/whispers/${whisperId}/comments`, { content })
}

export function listWhisperComments(whisperId: number) {
  return request.get<any, { id: number; content: string; created_at: string; author?: { nickname: string; avatar: string } }[]>(`/whispers/${whisperId}/comments`)
}
