import { Configuration, DefaultApi } from '@/client'
import { useAuthStore } from '@/store/auth'

export function useAPI() {
  const authStore = useAuthStore()
  return new DefaultApi(
    new Configuration({
      basePath: 'http://localhost:8090/api',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
        Authorization: `Bearer ${authStore.token}`
      }
    })
  )
}
