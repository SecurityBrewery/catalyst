import { defineStore } from 'pinia'
import { ref } from 'vue'

import type { User } from '@/client/models'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: ref<User | undefined>(undefined)
  }),
  getters: {
    isAuthenticated: (state) => !!state.user
  },
  actions: {
    setUser(user: User | undefined) {
      this.user = user
    }
  }
})
