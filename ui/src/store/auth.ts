import { defineStore } from 'pinia'
import { ref } from 'vue'

import type { User } from '@/client/models'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: ref<User | undefined>(undefined)
  }),
  actions: {
    setUser(user: User | undefined) {
      this.user = user
    }
  }
})
