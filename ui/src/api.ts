import { Configuration, DefaultApi } from '@/client'
import { useAuthStore } from '@/store/auth'

export function useAPI() {
  const authStore = useAuthStore()
  return new DefaultApi(
    new Configuration({
      basePath: '/api',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
        Authorization: `Bearer ${authStore.token}`
      }
    })
  )
}
