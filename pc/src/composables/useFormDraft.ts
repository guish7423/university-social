import { watch, onUnmounted } from "vue"
import type { Ref } from "vue"

const DRAFT_PREFIX = "draft_"

export function useFormDraft(key: string, fields: Record<string, Ref<string>>) {
  const storageKey = DRAFT_PREFIX + key

  // Load saved draft
  try {
    const saved = localStorage.getItem(storageKey)
    if (saved) {
      const data = JSON.parse(saved)
      for (const [k, ref] of Object.entries(fields)) {
        if (data[k] !== undefined) ref.value = data[k]
      }
    }
  } catch { /* ignore corrupt data */ }

  // Auto-save on changes (debounced via immediate watcher per field)
  const unwatch = watch(
    () => Object.fromEntries(Object.entries(fields).map(([k, r]) => [k, r.value])),
    (val) => {
      if (Object.values(val).some(v => v)) localStorage.setItem(storageKey, JSON.stringify(val))
    },
    { deep: true }
  )

  function clearDraft() {
    localStorage.removeItem(storageKey)
  }

  onUnmounted(unwatch)

  return { clearDraft }
}
