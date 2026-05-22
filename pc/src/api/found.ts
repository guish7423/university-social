import request from "./request"

export interface LostItemData {
  id: number
  user_id: number
  title: string
  description: string
  category: string
  location: string
  contact: string
  images: string[]
  status: number
  created_at: string
  is_owner?: boolean
}

export function listLostItems(params?: { category?: string; offset?: number; limit?: number }) {
  return request.get<any, LostItemData[]>("/lost-items", { params })
}

export function getLostItem(id: number) {
  return request.get<any, LostItemData>(`/lost-items/${id}`)
}

export function createLostItem(data: {
  title: string; description?: string; category: string; location?: string; contact?: string; images?: string[]
}) {
  return request.post<any, { id: number }>("/lost-items", data)
}

export function deleteLostItem(id: number) {
  return request.delete<any, { message: string }>(`/lost-items/${id}`)
}

export function updateLostItemStatus(id: number, status: number) {
  return request.put<any, { message: string }>(`/lost-items/${id}/status`, { status })
}

export function myLostItems() {
  return request.get<any, LostItemData[]>("/lost-items/my")
}

export function searchLostItems(q: string) {
  return request.get<any, LostItemData[]>("/lost-items/search", { params: { q } })
}
