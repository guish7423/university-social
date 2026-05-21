import { request } from "./request"

export function searchPosts(q: string, offset = 0, limit = 20) {
  return request<PostData[]>({ url: "/posts/search", method: "GET", data: { q, offset, limit } })
}

export function trendingPosts(days = 7, limit = 20) {
  return request<PostData[]>({ url: "/posts/trending", method: "GET", data: { days, limit } })
}

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

export interface UploadResult {
  url: string
}

export function uploadImage(tempFilePath: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync("token")
    uni.uploadFile({
      url: (import.meta.env.VITE_API_BASE || "http://localhost:8081/api/v1") + "/upload",
      filePath: tempFilePath,
      name: "file",
      header: { Authorization: "Bearer " + token },
      success: (res) => {
        if (res.statusCode === 200) {
          const data = JSON.parse(res.data as string)
          const base = (import.meta.env.VITE_API_BASE || "http://localhost:8081/api/v1").replace("/api/v1", "")
          resolve(base + data.url)
        } else {
          reject(new Error("上传失败"))
        }
      },
      fail: (err) => reject(err),
    })
  })
}

export function uploadImages(tempFilePaths: string[]): Promise<string[]> {
  return Promise.all(tempFilePaths.map(uploadImage))
}

export interface ShareCardData {
	title: string
	description: string
	image: string
	path: string
	author?: { nickname: string; avatar: string }
	like_count: number
}

export async function shareCard(postId: number): Promise<ShareCardData> {
	return request(`/posts/${postId}/share`)
}

export function followingPosts(offset = 0, limit = 20) {
	return request<PostData[]>({ url: "/posts/following", method: "GET", data: { offset, limit } })
}

export function userPosts(userId: number, offset = 0, limit = 20) {
	return request<PostData[]>({ url: `/users/${userId}/posts`, method: "GET", data: { offset, limit } })
}
