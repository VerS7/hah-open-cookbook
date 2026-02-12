import { ref, onUnmounted } from 'vue'

export function useDebounce<T>(initialValue: T, delay: number = 500) {
  const value = ref<T>(initialValue)
  const debouncedValue = ref<T>(initialValue)
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  const updateDebouncedValue = (newValue: T) => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    timeoutId = setTimeout(() => {
      debouncedValue.value = newValue
    }, delay)
  }

  onUnmounted(() => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
  })

  return {
    value,
    debouncedValue,
    updateDebouncedValue,
  }
}
