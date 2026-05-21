import { request } from "./request"

export function searchCircles(q: string) {
  return request<CircleData[]>({ url: "/circles/search", method: "GET", data: { q } })
}
import type { UserInfo } from "./user"

export interface CircleData {
  id: number
  name: string
  description: string
  icon: string
  cover: string
  creator_id: number
  member_count: number
  post_count: number
  created_at: string
  is_member: boolean
  creator?: { nickname: string; avatar: string }
}

export interface CircleMemberData {
  id: number
  circle_id: number
  user_id: number
  role: string
  created_at: string
  user: UserInfo
}

export interface CirclePostData {
  id: number
  circle_id: number
  user_id: number
  content: string
  image_urls: string[]
  like_count: number
  comment_count: number
  created_at: string
  author: UserInfo
  is_liked: boolean
}

export function listCircles(offset = 0, limit = 20) {
  return request<CircleData[]>({ url: "/circles", method: "GET", data: { offset, limit } })
}

export function getCircle(id: number) {
  return request<CircleData>({ url: `/circles/${id}`, method: "GET" })
}

export function createCircle(data: { name: string; description?: string; icon?: string; cover?: string }) {
  return request<{ id: number }>({ url: "/circles", method: "POST", data })
}

export function updateCircle(id: number, data: { name: string; description?: string; icon?: string; cover?: string }) {
  return request<{ message: string }>({ url: `/circles/${id}`, method: "PUT", data })
}

export function joinCircle(id: number) {
  return request<{ message: string }>({ url: `/circles/${id}/join`, method: "POST" })
}

export function leaveCircle(id: number) {
  return request<{ message: string }>({ url: `/circles/${id}/leave`, method: "POST" })
}

export function listCircleMembers(id: number) {
  return request<CircleMemberData[]>({ url: `/circles/${id}/members`, method: "GET" })
}

export function listCirclePosts(id: number, offset = 0, limit = 20) {
  return request<CirclePostData[]>({ url: `/circles/${id}/posts`, method: "GET", data: { offset, limit } })
}

export function createCirclePost(circleId: number, data: { content: string; image_urls?: string[] }) {
  return request<{ id: number }>({ url: `/circles/${circleId}/posts`, method: "POST", data })
}

export function toggleCirclePostLike(postId: number) {
  return request<{ liked: boolean }>({ url: `/circles/posts/${postId}/like`, method: "POST" })
}

export function createCirclePostComment(postId: number, content: string) {
  return request<{ id: number }>({ url: `/circles/posts/${postId}/comments`, method: "POST", data: { content } })
}

export function listCirclePostComments(postId: number) {
  return request<{ id: number; user_id: number; content: string; created_at: string; author: UserInfo }[]>({
    url: `/circles/posts/${postId}/comments`, method: "GET",
  })
}
