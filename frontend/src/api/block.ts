import { request } from "./request"
import type { UserInfo } from "./user"

export interface BlockStatus {
  is_blocked: boolean
}

export function blockUser(id: number) {
  return request<{ message: string }>({ url: "/users/" + id + "/block", method: "POST" })
}

export function unblockUser(id: number) {
  return request<{ message: string }>({ url: "/users/" + id + "/block", method: "DELETE" })
}

export function checkBlock(id: number) {
  return request<BlockStatus>({ url: "/users/" + id + "/block-status", method: "GET" })
}

export function getBlockedUsers() {
  return request<UserInfo[]>({ url: "/users/blocked", method: "GET" })
}
