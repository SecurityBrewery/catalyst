<script setup lang="ts">
import Toaster from '@/components/ui/toast/Toaster.vue'

import { onMounted, watch } from 'vue'
import { RouterView } from 'vue-router'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()

const router = useRouter()

const fetchUser = () => {
  if (!authStore.token) {
    authStore.setUser(undefined)
    authStore.setPermissions([])
    return
  }

  fetch('/auth/user', { headers: { Authorization: `Bearer ${authStore.token}` } }).then(
    (response) => {
      if (response.ok) {
        response.json().then((user) => {
          if (user) {
            authStore.setUser(user.user)
            authStore.setPermissions(user.permissions)
          } else {
            authStore.setUser(undefined)
            authStore.setPermissions([])
            router.push({ name: 'login' })
          }
        })
      }
    }
  )
}

onMounted(() => {
  fetchUser()
})

watch(
  () => authStore.token,
  () => fetchUser(),
  { immediate: true }
)
</script>

<template>
  <RouterView />
  <Toaster />
</template>
