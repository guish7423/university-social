import { request } from "./request"

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
  updated_at: string
  is_liked?: boolean
  is_owner?: boolean
}

export interface ProductComment {
  id: number
  product_id: number
  user_id: number
  content: string
  created_at: string
}

export function createProduct(data: {
  title: string
  description?: string
  price: number
  original_price?: number
  category: string
  condition?: string
  images?: string[]
  contact?: string
}) {
  return request<{ id: number }>({ url: "/products", method: "POST", data })
}

export function listProducts(params?: { category?: string; sort?: string; offset?: number; limit?: number }) {
  return request<ProductData[]>({ url: "/products", method: "GET", data: params as any })
}

export function getProduct(id: number) {
  return request<ProductData>({ url: `/products/${id}`, method: "GET" })
}

export function deleteProduct(id: number) {
  return request<{ message: string }>({ url: `/products/${id}`, method: "DELETE" })
}

export function markSold(id: number) {
  return request<{ message: string }>({ url: `/products/${id}/sold`, method: "PUT" })
}

export function myProducts(params?: { offset?: number; limit?: number }) {
  return request<{ products: ProductData[]; total: number; active: number }>({ url: "/products/my", method: "GET", data: params as any })
}

export function searchProducts(q: string, offset = 0, limit = 20) {
  return request<ProductData[]>({ url: "/products/search", method: "GET", data: { q, offset, limit } })
}

export function toggleProductLike(id: number) {
  return request<{ liked: boolean }>({ url: `/products/${id}/like`, method: "POST" })
}

export function createProductComment(productId: number, content: string) {
  return request<{ id: number }>({ url: `/products/${productId}/comments`, method: "POST", data: { content } })
}

export function listProductComments(productId: number) {
  return request<ProductComment[]>({ url: `/products/${productId}/comments`, method: "GET" })
}

export function productCategories() {
  return request<string[]>({ url: "/products/categories", method: "GET" })
}

export function trendingProducts() {
  return request<ProductData[]>({ url: "/products/trending", method: "GET" })
}
