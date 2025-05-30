import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'

export const useCatalystStore = defineStore('catalyst', {
  state: () => ({
    sidebarCollapsed: useLocalStorage('sidebarCollapsed', false)
  }),
  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    }
  }
})
