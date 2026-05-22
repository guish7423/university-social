import { request } from "./request"
import type { UserInfo } from "./user"
import type { PostData } from "./post"

export interface SearchUserResult extends UserInfo {
  friend_status: string
}

export function searchUsers(q: string) {
  return request<SearchUserResult[]>({ url: "/users/search", method: "GET", data: { q } })
}
export function listFriends() {
  return request<UserInfo[]>({ url: "/friends", method: "GET" })
}

export function listFriendRequests() {
  return request<UserInfo[]>({ url: "/friends/requests", method: "GET" })
}

export function sendFriendRequest(id: number) {
  return request<{ message: string }>({ url: `/friends/request/${id}`, method: "POST" })
}

export function acceptFriendRequest(id: number) {
  return request<{ message: string }>({ url: `/friends/accept/${id}`, method: "POST" })
}
export function rejectFriendRequest(id: number) {
  return request<{ message: string }>({ url: `/friends/reject/${id}`, method: "POST" })
}

export function removeFriend(id: number) {
  return request<{ message: string }>({ url: `/friends/${id}`, method: "DELETE" })
}

export interface NotificationData {
  id: number
  user_id: number
  from_user_id: number
  type: string
  content: string
  ref_id: number
  is_read: boolean
  created_at: string
  from_user?: { nickname: string; avatar: string }
}

export function listNotifications() {
  return request<NotificationData[]>({ url: "/notifications", method: "GET" })
}

export function markNotificationsRead() {
  return request<{ message: string }>({ url: "/notifications/read", method: "POST" })
}


export interface BannerData {
  id: number
  title: string
  image_url: string
  link_url: string
  sort_order: number
  is_active: boolean
  created_at: string
}

export interface AnnouncementData {
  id: number
  title: string
  content: string
  is_active: boolean
  created_at: string
}

export function listBanners() {
  return request<BannerData[]>({ url: "/banners", method: "GET" })
}

export function listAnnouncements() {
  return request<AnnouncementData[]>({ url: "/announcements", method: "GET" })
}

export function listFeaturedPosts(offset = 0, limit = 20) {
  return request<PostData[]>({ url: "/posts/featured", method: "GET", data: { offset, limit } })
}