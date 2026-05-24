import request from "./request"

export interface PostData {
  id: number
  user_id: number
  title?: string
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
  return request.post<any, { id: number }>("/posts", data)
}

export function listPosts(params?: { school?: string; topic_id?: number; offset?: number; limit?: number }) {
  return request.get<any, PostData[]>("/posts", { params })
}

export function getPost(id: number) {
  return request.get<any, PostData>(`/posts/${id}`)
}

export function deletePost(id: number) {
  return request.delete<any, { message: string }>(`/posts/${id}`)
}

export function searchPosts(q: string, offset = 0, limit = 20) {
  return request.get<any, PostData[]>("/posts/search", { params: { q, offset, limit } })
}

export function trendingPosts(days = 7, limit = 20) {
  return request.get<any, PostData[]>("/posts/trending", { params: { days, limit } })
}

export function followingPosts(offset = 0, limit = 20) {
  return request.get<any, PostData[]>("/posts/following", { params: { offset, limit } })
}

export function userPosts(userId: number, offset = 0, limit = 20) {
  return request.get<any, PostData[]>(`/users/${userId}/posts`, { params: { offset, limit } })
}

export function createComment(postId: number, content: string) {
  return request.post<any, { id: number }>(`/posts/${postId}/comments`, { content })
}

export function listComments(postId: number, offset = 0, limit = 20) {
  return request.get<any, CommentData[]>(`/posts/${postId}/comments`, { params: { offset, limit } })
}

export function toggleLike(postId: number) {
  return request.post<any, { liked: boolean }>(`/posts/${postId}/like`)
}

export function listTopics() {
  return request.get<any, TopicData[]>("/topics")
}
