import { request } from "./request"

export interface Course {
  id: number
  name: string
  teacher: string
  school: string
  department: string
  created_at: string
}

export interface CourseRating {
  id: number
  course_id: number
  user_id: number
  teaching_rating: number
  difficulty: number
  grading: number
  comment: string
  is_anonymous: boolean
  helpful_count: number
  created_at: string
  nickname?: string
  avatar?: string
}

export interface CourseDetail {
  id: number
  name: string
  teacher: string
  school: string
  department: string
  avg_teaching: number
  avg_difficulty: number
  avg_grading: number
  rating_count: number
  user_rating?: CourseRating | null
  created_at: string
}

export function searchCourses(q: string) {
  return request<Course[]>({ url: "/courses/search", method: "GET", data: { q } as any })
}

export function getCourseDetail(id: number) {
  return request<CourseDetail>({ url: `/courses/${id}`, method: "GET" })
}

export function createRating(courseId: number, data: {
  teaching_rating: number
  difficulty: number
  grading: number
  comment?: string
  is_anonymous?: boolean
}) {
  return request<{ id: number }>({ url: `/courses/${courseId}/rating`, method: "POST", data })
}

export function listRatings(courseId: number, params?: { order?: string; offset?: number; limit?: number }) {
  return request<CourseRating[]>({ url: `/courses/${courseId}/ratings`, method: "GET", data: params as any })
}

export function getUserRating(courseId: number) {
  return request<CourseRating | null>({ url: `/courses/${courseId}/my-rating`, method: "GET" })
}

export function markHelpful(ratingId: number) {
  return request<{ message: string }>({ url: `/courses/rating/${ratingId}/helpful`, method: "POST" })
}
