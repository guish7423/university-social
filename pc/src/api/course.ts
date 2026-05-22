import request from "./request"

export interface CourseData {
  id: number
  name: string
  teacher: string
  department: string
  week: string
  time: string
  location: string
  capacity: number
  enrolled: number
  credits: number
  rating: number
}

export function searchCourses(q: string) {
  return request.get<any, CourseData[]>("/courses/search", { params: { q } })
}

export function getCourseDetail(id: number) {
  return request.get<any, CourseData>(`/courses/${id}`)
}

export function createRating(courseId: number, data: { rating: number; comment?: string }) {
  return request.post<any, { id: number }>(`/courses/${courseId}/rating`, data)
}

export function listRatings(courseId: number) {
  return request.get<any, { id: number; rating: number; comment: string; created_at: string; author?: { nickname: string } }[]>(`/courses/${courseId}/ratings`)
}
