import request from "./request"

export function followUser(id: number) {
  return request.post<any, { message: string }>(`/follow/${id}`)
}

export function unfollowUser(id: number) {
  return request.delete<any, { message: string }>(`/follow/${id}`)
}

export function isFollowing(id: number) {
  return request.get<any, { following: boolean }>(`/follow/${id}/check`)
}

export function followingList(id: number) {
  return request.get<any, { id: number; nickname: string; avatar: string }[]>(`/users/${id}/following`)
}

export function followersList(id: number) {
  return request.get<any, { id: number; nickname: string; avatar: string }[]>(`/users/${id}/followers`)
}

export function followCounts(id: number) {
  return request.get<any, { following: number; followers: number }>(`/users/${id}/follow-counts`)
}
