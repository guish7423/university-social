import request from "./request"

export function addFavorite(postId: number) {
  return request.post<any, { message: string }>(`/posts/${postId}/favorite`)
}

export function removeFavorite(postId: number) {
  return request.delete<any, { message: string }>(`/posts/${postId}/favorite`)
}

export function getFavoriteStatus(postId: number) {
  return request.get<any, { is_favorited: boolean }>(`/posts/${postId}/favorite-status`)
}

export function listFavorites(page = 1, limit = 20) {
  return request.get<any, { posts: any[]; total: number; page: number }>("/favorites", {
    params: { page, limit },
  })
}
