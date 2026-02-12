import { ref, type Ref } from 'vue'

export const BASE_URL = import.meta.env.VITE_API_URL || ''
export const BASE_API_URL = BASE_URL + '/api'

export interface Header {
  [key: string]: string
}

export interface APIResponse<T> {
  data: Ref<T | null>
  error: Ref<string | null>
  loading: Ref<boolean>
  execute: () => Promise<void>
}

export function useAPI<T>(url: string, request?: RequestInit): APIResponse<T> {
  const data = ref<T | null>(null) as Ref<T | null>
  const error = ref<string | null>(null)
  const loading = ref(false)

  async function execute(): Promise<void> {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(url, request)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      if (response.headers.get('content-type')?.includes('application/json')) {
        data.value = await response.json()
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'uncommon error'
    } finally {
      loading.value = false
    }
  }

  return {
    data,
    error,
    loading,
    execute,
  }
}
