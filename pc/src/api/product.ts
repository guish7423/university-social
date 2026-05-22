import request from "./request"

export interface ProductData {
  id: number
  user_id: number
  title: string
  description: string
  price: number
  original_price?: number
  category: string
  condition: string
  images: string[]
  contact: string
  status: number
  like_count: number
  comment_count: number
  created_at: string
  is_liked?: boolean
  is_owner?: boolean
}

export function listProducts(params?: { category?: string; sort?: string; offset?: number; limit?: number }) {
  return request.get<any, ProductData[]>("/products", { params })
}

export function getProduct(id: number) {
  return request.get<any, ProductData>(`/products/${id}`)
}

export function createProduct(data: {
  title: string; description?: string; price: number; original_price?: number
  category: string; condition?: string; images?: string[]; contact?: string
}) {
  return request.post<any, { id: number }>("/products", data)
}

export function deleteProduct(id: number) {
  return request.delete<any, { message: string }>(`/products/${id}`)
}

export function markSold(id: number) {
  return request.put<any, { message: string }>(`/products/${id}/sold`)
}

export function myProducts(params?: { offset?: number; limit?: number }) {
  return request.get<any, { products: ProductData[]; total: number; active: number }>("/products/my", { params })
}

export function toggleProductLike(id: number) {
  return request.post<any, { liked: boolean }>(`/products/${id}/like`)
}

export function productCategories() {
  return request.get<any, string[]>("/products/categories")
}

export function searchProducts(q: string, offset = 0, limit = 20) {
  return request.get<any, ProductData[]>("/products/search", { params: { q, offset, limit } })
}
