import { request } from './request'

export interface LostItem {
  id: number
  user_id: number
  title: string
  description: string
  category: string
  item_type: string
  location: string
  contact: string
  images: string[]
  status: number
  created_at: string
  updated_at: string
  nickname: string
  avatar: string
}

export interface CreateLostItemData {
  title: string
  description: string
  category: string
  item_type?: string
  location?: string
  contact?: string
  images?: string[]
}

export function listLostItems(category: string = '', status: number = -1, offset: number = 0, limit: number = 20) {
  return request.get<LostItem[]>('/lost-items', { params: { category, status, offset, limit } })
}

export function createLostItem(data: CreateLostItemData) {
  return request.post<LostItem>('/lost-items', data)
}

export function getLostItem(id: number) {
  return request.get<LostItem>(`/lost-items/${id}`)
}

export function deleteLostItem(id: number) {
  return request.delete(`/lost-items/${id}`)
}

export function updateLostItemStatus(id: number, status: number) {
  return request.put(`/lost-items/${id}/status`, { status })
}

export function myLostItems(offset: number = 0, limit: number = 20) {
  return request.get<LostItem[]>('/lost-items/my', { params: { offset, limit } })
}

export function searchLostItems(q: string, category: string = '', offset: number = 0, limit: number = 20) {
  return request.get<LostItem[]>('/lost-items/search', { params: { q, category, offset, limit } })
}
