import request from "./request"
import type { UserInfo } from "./auth"
import type { PostData } from "./post"

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

export interface BannerData {
  id: number
  title: string
  image_url: string
  link_url: string
  sort_order: number
  is_active: boolean
}

export interface AnnouncementData {
  id: number
  title: string
  content: string
  is_active: boolean
  created_at: string
}

export function searchUsers(q: string) {
  return request.get<any, (UserInfo & { friend_status: string })[]>("/users/search", { params: { q } })
}

export function listFriends() {
  return request.get<any, UserInfo[]>("/friends")
}

export function listFriendRequests() {
  return request.get<any, UserInfo[]>("/friends/requests")
}

export function sendFriendRequest(id: number) {
  return request.post<any, { message: string }>(`/friends/request/${id}`)
}

export function acceptFriendRequest(id: number) {
  return request.post<any, { message: string }>(`/friends/accept/${id}`)
}

export function rejectFriendRequest(id: number) {
  return request.post<any, { message: string }>(`/friends/reject/${id}`)
}

export function removeFriend(id: number) {
  return request.delete<any, { message: string }>(`/friends/${id}`)
}

export function listNotifications() {
  return request.get<any, NotificationData[]>("/notifications")
}

export function markNotificationsRead() {
  return request.post<any, { message: string }>("/notifications/read")
}

export function listBanners() {
  return request.get<any, BannerData[]>("/banners")
}

export function listAnnouncements() {
  return request.get<any, AnnouncementData[]>("/announcements")
}

export function listFeaturedPosts(offset = 0, limit = 20) {
  return request.get<any, PostData[]>("/posts/featured", { params: { offset, limit } })
}

export function getFollowers(userId: number) {
  return request.get<any, UserInfo[]>(`/users/${userId}/followers`)
}

export function getFollowing(userId: number) {
  return request.get<any, UserInfo[]>(`/users/${userId}/following`)
}

export function followUser(id: number) {
  return request.post<any, { message: string }>(`/follow/${id}`)
}

export function unfollowUser(id: number) {
  return request.delete<any, { message: string }>(`/follow/${id}`)
}

export function isFollowing(id: number) {
  return request.get<any, { following: boolean }>(`/follow/${id}/check`)
}

export function getFollowCounts(userId: number) {
  return request.get<any, { following: number; followers: number }>(`/users/${userId}/follow-counts`)
}
