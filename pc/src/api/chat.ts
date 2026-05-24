import request from "./request"

export interface ConversationData {
  user_id: number
  nickname: string
  avatar: string
  last_message: string
  last_time: string
  unread: number
}

export interface MessageData {
  id: number
  sender_id: number
  receiver_id: number
  content: string
  msg_type: number
  created_at: string
  read_at: string | null
  sender?: { nickname: string; avatar: string }
}

export function listConversations() {
  return request.get<any, ConversationData[]>("/chat/conversations")
}

export function getMessages(userId: number) {
  return request.get<any, MessageData[]>(`/chat/messages/${userId}`)
}

export function sendMessage(toUserId: number, content: string) {
  return request.post<any, { id: number }>("/chat/send", { to_user_id: toUserId, content })
}

export function markChatRead(userId: number) {
  return request.post<any, { message: string }>(`/chat/read/${userId}`)
}

export function unreadCount() {
  return request.get<any, { count: number }>("/chat/unread")
}
