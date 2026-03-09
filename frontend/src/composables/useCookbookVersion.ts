import { computed, ref } from 'vue'
import { BASE_API_URL, useAPI } from './useAPI'

export interface CookbookVersion {
  version: string
  isArchived: boolean
}

export function useCookbookVersions() {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const rawVersions = ref<string[] | null>(null)

  const versions = computed<CookbookVersion[] | null>(() => {
    return rawVersions.value !== null
      ? rawVersions.value.map(
          (v) => ({ version: v, isArchived: v.includes('archived') }) as CookbookVersion,
        )
      : null
  })

  async function get() {
    loading.value = true
    error.value = null

    const {
      data,
      error: APIerror,
      execute,
    } = useAPI<string[]>(BASE_API_URL + `/versions`, {
      method: 'GET',
    })

    await execute()
    if (APIerror.value) error.value = APIerror.value
    if (data.value) rawVersions.value = data.value
    loading.value = false
  }

  return { versions, loading, error, get }
}
