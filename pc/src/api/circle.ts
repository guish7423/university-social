import request from "./request"

export interface CircleData {
  id: number
  name: string
  description: string
  avatar: string
  member_count: number
  post_count: number
  join_type: string
  has_pending_request: boolean
  is_joined: boolean
  creator_id: number
  created_at: string
}

export interface JoinRequestData {
  id: number
  circle_id: number
  user_id: number
  status: string
  created_at: string
  handled_at?: string
  handler_id?: number
  user?: { nickname: string; avatar: string }
}

export interface CirclePostData {
  id: number
  circle_id: number
  user_id: number
  content: string
  images: string[]
  like_count: number
  comment_count: number
  created_at: string
  author?: { nickname: string; avatar: string }
  is_liked?: boolean
}

export function listCircles(params?: { offset?: number; limit?: number }) {
  return request.get<any, CircleData[]>("/circles", { params })
}

export function getCircle(id: number) {
  return request.get<any, CircleData>(`/circles/${id}`)
}

export function createCircle(data: { name: string; description?: string; avatar?: string }) {
  return request.post<any, { id: number }>("/circles", data)
}

export function updateCircle(id: number, data: { name?: string; description?: string; avatar?: string }) {
  return request.put<any, { message: string }>(`/circles/${id}`, data)
}

export function joinCircle(id: number) {
  return request.post<any, { message: string }>(`/circles/${id}/join`)
}

export function leaveCircle(id: number) {
  return request.post<any, { message: string }>(`/circles/${id}/leave`)
}

export function listCircleMembers(id: number) {
  return request.get<any, { id: number; nickname: string; avatar: string }[]>(`/circles/${id}/members`)
}

export function listCirclePosts(id: number, offset = 0, limit = 20) {
  return request.get<any, CirclePostData[]>(`/circles/${id}/posts`, { params: { offset, limit } })
}

export function createCirclePost(circleId: number, data: { content: string; images?: string[] }) {
  return request.post<any, { id: number }>(`/circles/${circleId}/posts`, data)
}

export function searchCircles(q: string) {
  return request.get<any, CircleData[]>("/circles/search", { params: { q } })
}


export function listJoinRequests(circleId: number) {
  return request.get<any, JoinRequestData[]>(`/circles/${circleId}/join-requests`)
}

export function handleJoinRequest(circleId: number, requestId: number, action: 'approve' | 'reject') {
  return request.post<any, { message: string }>(`/circles/join-requests/${requestId}/handle`, { action })
}
