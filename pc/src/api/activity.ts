import request from "./request"

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

export function listActivities(params?: { type?: string; offset?: number; limit?: number }) {
  return request.get<any, ActivityData[]>("/activities", { params })
}

export function getActivity(id: number) {
  return request.get<any, ActivityData>(`/activities/${id}`)
}

export function createActivity(data: {
  title: string; description?: string; activity_type: string; location?: string
  max_participants?: number; start_time: string; end_time?: string; images?: string[]
}) {
  return request.post<any, { id: number }>("/activities", data)
}

export function deleteActivity(id: number) {
  return request.delete<any, { message: string }>(`/activities/${id}`)
}

export function joinActivity(id: number) {
  return request.post<any, { message: string }>(`/activities/${id}/join`)
}

export function leaveActivity(id: number) {
  return request.post<any, { message: string }>(`/activities/${id}/leave`)
}

export function activityTypes() {
  return request.get<any, string[]>("/activities/types")
}

export function myActivities() {
  return request.get<any, { created: ActivityData[]; joined: ActivityData[] }>("/activities/my")
}
