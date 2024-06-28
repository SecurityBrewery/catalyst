import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import type { Ref } from 'vue'

type State = {
  sidebarCollapsed: Ref<boolean>
}

export const useCatalystStore = defineStore('catalyst', {
  state: (): State => ({
    sidebarCollapsed: useLocalStorage('sidebarCollapsed', false)
  }),
  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    }
  }
})
