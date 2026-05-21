import { request } from "./request"

export interface ActivityData {
  id: number
  user_id: number
  title: string
  description: string
  activity_type: string
  location: string
  max_participants: number
  participant_count: number
  comment_count: number
  start_time: string
  end_time?: string
  images: string[]
  status: number
  created_at: string
  updated_at: string
  is_participant?: boolean
  is_owner?: boolean
}

export interface ActivityComment {
  id: number
  activity_id: number
  user_id: number
  content: string
  created_at: string
}

export function createActivity(data: {
  title: string
  description?: string
  activity_type: string
  location?: string
  max_participants?: number
  start_time: string
  end_time?: string
  images?: string[]
}) {
  return request<{ id: number }>({ url: "/activities", method: "POST", data })
}

export function listActivities(params?: { type?: string; offset?: number; limit?: number }) {
  return request<ActivityData[]>({ url: "/activities", method: "GET", data: params as any })
}

export function getActivity(id: number) {
  return request<ActivityData>({ url: `/activities/${id}`, method: "GET" })
}

export function deleteActivity(id: number) {
  return request<{ message: string }>({ url: `/activities/${id}`, method: "DELETE" })
}

export function joinActivity(id: number) {
  return request<{ message: string }>({ url: `/activities/${id}/join`, method: "POST" })
}

export function leaveActivity(id: number) {
  return request<{ message: string }>({ url: `/activities/${id}/leave`, method: "POST" })
}

export function activityParticipants(id: number) {
  return request<{ id: number; nickname: string; avatar: string }[]>({ url: `/activities/${id}/participants`, method: "GET" })
}

export function createActivityComment(activityId: number, content: string) {
  return request<{ id: number }>({ url: `/activities/${activityId}/comments`, method: "POST", data: { content } })
}

export function listActivityComments(activityId: number) {
  return request<ActivityComment[]>({ url: `/activities/${activityId}/comments`, method: "GET" })
}

export function myActivities() {
  return request<{ created: ActivityData[]; joined: ActivityData[] }>({ url: "/activities/my", method: "GET" })
}

export function activityTypes() {
  return request<string[]>({ url: "/activities/types", method: "GET" })
}
