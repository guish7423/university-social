import { ref, watch, onUnmounted } from "vue"

export function useAnimatedNumber(source: import("vue").Ref<number>, duration = 300) {
  const display = ref(source.value)
  let rafId = 0

  watch(source, (newVal) => {
    cancelAnimationFrame(rafId)
    const start = display.value
    const diff = newVal - start
    if (diff === 0) return

    const startTime = performance.now()

    function animate(currentTime: number) {
      const elapsed = currentTime - startTime
      const progress = Math.min(elapsed / duration, 1)
      const eased = 1 - Math.pow(1 - progress, 3) // ease-out cubic
      display.value = Math.round(start + diff * eased)

      if (progress < 1) {
        rafId = requestAnimationFrame(animate)
      }
    }

    rafId = requestAnimationFrame(animate)
  })

  onUnmounted(() => {
    cancelAnimationFrame(rafId)
  })

  return display
}
