import { request } from "./request"

export interface PostData {
  id: number
  user_id: number
  content: string
  images: string[]
  topic_id?: number
  school: string
  like_count: number
  comment_count: number
  created_at: string
  updated_at: string
  author?: { nickname: string; avatar: string }
  is_liked?: boolean
}

export interface CommentData {
  id: number
  post_id: number
  user_id: number
  content: string
  created_at: string
  author?: { nickname: string; avatar: string }
}

export interface TopicData {
  id: number
  name: string
  icon: string
  post_count: number
}

export function createPost(data: { content: string; images?: string[]; topic_id?: number }) {
  return request<{ id: number }>({ url: "/posts", method: "POST", data })
}

export function listPosts(params?: { school?: string; topic_id?: number; offset?: number; limit?: number }) {
  return request<PostData[]>({ url: "/posts", method: "GET", data: params as any })
}

export function getPost(id: number) {
  return request<PostData>({ url: `/posts/${id}`, method: "GET" })
}

export function deletePost(id: number) {
  return request<{ message: string }>({ url: `/posts/${id}`, method: "DELETE" })
}

export function createComment(postId: number, content: string) {
  return request<{ id: number }>({ url: `/posts/${postId}/comments`, method: "POST", data: { content } })
}

export function listComments(postId: number, offset = 0, limit = 20) {
  return request<CommentData[]>({ url: `/posts/${postId}/comments`, method: "GET", data: { offset, limit } })
}

export function toggleLike(postId: number) {
  return request<{ liked: boolean }>({ url: `/posts/${postId}/like`, method: "POST" })
}

export function listTopics() {
  return request<TopicData[]>({ url: "/topics", method: "GET" })
}
