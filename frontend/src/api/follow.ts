import { request } from "./request"

export interface FollowUser {
  id: number
  nickname: string
  avatar: string
  school: string
  is_verified: boolean
}

export interface FollowItem {
  id: number
  follower_id: number
  following_id: number
  created_at: string
  following?: FollowUser
  follower?: FollowUser
}

export function followUser(id: number) {
  return request<{ message: string }>({ url: "/follow/" + id, method: "POST" })
}

export function unfollowUser(id: number) {
  return request<{ message: string }>({ url: "/follow/" + id, method: "DELETE" })
}

export function checkFollow(id: number) {
  return request<{ is_following: boolean }>({ url: "/follow/" + id + "/check", method: "GET" })
}

export async function getFollowing(userId: number) {
  const data = await request<FollowItem[]>({ url: "/users/" + userId + "/following", method: "GET" })
  return data.map(item => item.following!).filter(Boolean)
}

export async function getFollowers(userId: number) {
  const data = await request<FollowItem[]>({ url: "/users/" + userId + "/followers", method: "GET" })
  return data.map(item => item.follower!).filter(Boolean)
}

export async function getFollowCounts(userId: number) {
  const data = await request<{ followers: number; following: number }>({ url: "/users/" + userId + "/follow-counts", method: "GET" })
  return { followers_count: data.followers, following_count: data.following }
}
