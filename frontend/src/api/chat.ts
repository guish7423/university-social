import { request } from "./request"

export interface MessageData {
  id: number
  sender_id: number
  receiver_id: number
  content: string
  msg_type: number
  created_at: string
  read_at?: string
  sender?: { nickname: string; avatar: string }
}

export interface ConversationData {
  other_user_id: number
  other_nickname: string
  other_avatar: string
  last_content: string
  last_time: string
  unread_count: number
}

export function sendMessage(data: { receiver_id: number; content: string; msg_type?: number }) {
  return request<MessageData>({ url: "/chat/send", method: "POST", data })
}

export function getConversations() {
  return request<ConversationData[]>({ url: "/chat/conversations", method: "GET" })
}

export function getMessages(userId: number, params?: { offset?: number; limit?: number }) {
  return request<MessageData[]>({ url: `/chat/messages/${userId}`, method: "GET", data: params as any })
}

export function markRead(userId: number) {
  return request<{ ok: boolean }>({ url: `/chat/read/${userId}`, method: "POST" })
}

export function getUnreadCount() {
  return request<{ count: number }>({ url: "/chat/unread", method: "GET" })
}
