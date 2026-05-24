import request from "./request"

export interface ReviewTargetData {
  id: number
  name: string
  category: string
  description: string
  creatorId: number
  creator: { id: number; nickname: string; avatar?: string }
  avgScore: number
  totalRatings: number
  distribution: number[]
  createdAt: string
}

export interface RatingData {
  id: number
  targetId: number
  userId: number
  user: { id: number; nickname: string; avatar?: string }
  score: number
  tags: string[]
  content: string
  helpfulCount: number
  isHelpfulByMe: boolean
  createdAt: string
}

export function listReviewTargets(params?: {
  page?: number; limit?: number
  category?: string; search?: string; sort?: string
}) {
  return request.get<any, { targets: ReviewTargetData[]; total: number }>("/review-targets", { params })
}

export function getReviewTarget(id: number) {
  return request.get<any, { target: ReviewTargetData; stats: { avgScore: number; distribution: number[]; totalRatings: number } }>(`/review-targets/${id}`)
}

export function createReviewTarget(data: { name: string; category: string; description?: string }) {
  return request.post<any, ReviewTargetData>("/review-targets", data)
}

export function searchReviewTargets(q: string) {
  return request.get<any, { targets: ReviewTargetData[] }>("/review-targets/search", { params: { q } })
}

export function listRatings(targetId: number, params?: { page?: number; limit?: number; sort?: string }) {
  return request.get<any, { ratings: RatingData[]; total: number }>(`/review-targets/${targetId}/ratings`, { params })
}

export function createRating(targetId: number, data: { score: number; tags?: string[]; content?: string }) {
  return request.post<any, RatingData>(`/review-targets/${targetId}/ratings`, data)
}

export function toggleHelpful(ratingId: number) {
  return request.post<any, { helpfulCount: number; isHelpfulByMe: boolean }>(`/ratings/${ratingId}/helpful`)
}
