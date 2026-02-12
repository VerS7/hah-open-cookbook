// useAuth.ts
import { ref, computed } from 'vue'
import { type Header, BASE_API_URL, useAPI } from './useAPI'

export interface UserAccount {
  username: string
  token: string
  isAdmin: boolean
}

const STORAGE_KEY = 'USER_ACCOUNT'

export function buildAuthTokenHeader(token: string): Header {
  return { Authorization: `Bearer ${token}` }
}

export function useAuth() {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const user = ref<UserAccount | null>(null)

  try {
    loadFromStorage()
  } catch {}

  function loadFromStorage() {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        user.value = JSON.parse(stored)
      }
    } catch (e) {
      clearStorage()
      throw e
    }
  }

  function saveToStorage(userData: UserAccount) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(userData))
  }

  function clearStorage() {
    localStorage.removeItem(STORAGE_KEY)
  }

  const isAdmin = computed<boolean | null>(() => (user.value ? user.value.isAdmin : null))
  const token = computed<string | null>(() => (user.value ? user.value.token : null))
  const username = computed<string | null>(() => (user.value ? user.value.username : null))
  const authHeader = computed<Header | null>(() =>
    user.value ? buildAuthTokenHeader(user.value.token) : null,
  )
  const isAuthenticated = computed(() => !!user.value)

  async function login(username: string, password: string) {
    loading.value = true
    error.value = null

    const query = new URLSearchParams({
      username: username,
      password: password,
    })

    const {
      data,
      error: APIerror,
      execute,
    } = useAPI<UserAccount>(BASE_API_URL + '/login', {
      method: 'POST',
      body: query,
    })

    try {
      await execute()

      if (APIerror.value) {
        error.value = APIerror.value
        return
      }

      if (data.value) {
        user.value = data.value
        saveToStorage(data.value)
        return
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
      return
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    clearStorage()
  }

  function checkAuth(): boolean {
    return isAuthenticated.value
  }

  return {
    user,
    loading,
    error,

    isAdmin,
    token,
    username,
    authHeader,
    isAuthenticated,

    login,
    logout,
    checkAuth,
    loadFromStorage,
  }
}
