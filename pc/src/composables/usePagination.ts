import { ref, type Ref } from "vue"

interface PaginationOptions<T> {
  fetchFn: (offset: number, limit: number) => Promise<T[]>
  pageSize?: number
  append?: boolean
}

export function usePagination<T>({ fetchFn, pageSize = 20, append = true }: PaginationOptions<T>) {
  const items = ref<T[]>([]) as Ref<T[]>
  const loading = ref(false)
  const hasMore = ref(true)
  const offset = ref(0)

  async function loadMore() {
    if (loading.value || !hasMore.value) return
    loading.value = true
    try {
      const res = await fetchFn(offset.value, pageSize)
      if (res.length < pageSize) hasMore.value = false
      if (append) {
        items.value = [...items.value, ...res] as T[]
      } else {
        items.value = res as T[]
      }
      offset.value += res.length
    } catch {
      // handled by interceptor
    } finally {
      loading.value = false
    }
  }

  async function reset() {
    items.value = [] as T[]
    offset.value = 0
    hasMore.value = true
    await loadMore()
  }

  return { items, loading, hasMore, loadMore, reset }
}
