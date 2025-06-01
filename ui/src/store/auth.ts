import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import { ref } from 'vue'

import type { User } from '@/client/models'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: useLocalStorage('token', ''),
    user: ref<User | undefined>(undefined)
  }),
  getters: {
    isAuthenticated: (state) => !!state.token
  },
  actions: {
    setToken(token: string) {
      this.token = token
    },
    setUser(user: User | undefined) {
      this.user = user
    }
  }
})
