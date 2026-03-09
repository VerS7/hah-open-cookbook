import { ref } from 'vue'
import { BASE_API_URL } from './useAPI'
import { buildAuthTokenHeader } from './useAuth'

export type Export = 'default' | 'nurgling' | 'json'

export function useExport(token: string, cookbookVersion: string) {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const exported = ref()

  async function exportAs(exportType: Export) {
    const query = new URLSearchParams({
      type: exportType,
    })

    loading.value = true
    error.value = null

    const response = await fetch(BASE_API_URL + `/${cookbookVersion}/export?${query}`, {
      method: 'GET',
      headers: buildAuthTokenHeader(token),
    })

    const url = window.URL.createObjectURL(await response.blob())
    const link = document.createElement('a')
    link.href = url

    switch (exportType) {
      case 'default':
        link.download = 'cookbook.db'
        break
      case 'nurgling':
        link.download = 'cookbook.db'
        break
      case 'json':
        link.download = 'cookbook.json'
        break
    }

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    loading.value = false
  }

  return { exported, loading, error, exportAs }
}
