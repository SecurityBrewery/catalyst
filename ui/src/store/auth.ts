import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import { ref } from 'vue'

import type { User } from '@/client/models'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: useLocalStorage('token', ''),
    user: ref<User | undefined>(undefined),
    permissions: ref<string[]>([])
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    hasPermission: (state) => (permission: string) => {
      return state.permissions.includes(permission) || state.permissions.includes('admin')
    }
  },
  actions: {
    setToken(token: string) {
      this.token = token
    },
    setUser(user: User | undefined) {
      this.user = user
    },
    setPermissions(permissions: string[]) {
      if (!permissions) {
        permissions = []
      }
      this.permissions = permissions
    }
  }
})
