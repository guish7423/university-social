import { request } from "./request"

export interface FavoriteData {
  id: number
  user_id: number
  post_id: number
  created_at: string
  post?: {
    id: number
    content: string
    images: string[]
    like_count: number
    comment_count: number
    created_at: string
    author?: { nickname: string; avatar: string }
  }
}

export function listFavorites(offset = 0, limit = 20) {
  return request<FavoriteData[]>({ url: "/favorites", method: "GET", data: { offset, limit } })
}

export function addFavorite(postId: number) {
  return request<{ message: string }>({ url: `/favorites/${postId}`, method: "POST" })
}

export function removeFavorite(postId: number) {
  return request<{ message: string }>({ url: `/favorites/${postId}`, method: "DELETE" })
}
