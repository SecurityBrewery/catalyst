<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ResourceListElement from '@/components/layout/ResourceListElement.vue'
import { Button } from '@/components/ui/button'

import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { User } from '@/client/models'

const api = useAPI()

const route = useRoute()
const router = useRouter()

const {
  isPending,
  isError,
  data: users,
  error
} = useQuery({
  queryKey: ['users'],
  queryFn: (): Promise<Array<User>> => api.listUsers()
})

const description = (user: User): string => {
  var desc = user.email

  if (!user.active) {
    desc += ' (inactive)'
  }

  return desc
}

const openNew = () => {
  router.push({ name: 'users', params: { id: 'new' } })
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader title="Users">
      <div class="ml-auto">
        <Button variant="ghost" @click="openNew">New User</Button>
      </div>
    </ColumnHeader>
    <div class="mt-2 flex flex-1 flex-col gap-2 overflow-scroll p-2 pt-0">
      <ResourceListElement
        v-for="user in users"
        :key="user.id"
        :title="user.name ? user.name : user.username"
        :created="user.created"
        :subtitle="user.username"
        :description="description(user)"
        :active="route.params.id === user.id"
        :to="{ name: 'users', params: { id: user.id } }"
        :open="false"
      >
        {{ user.name }}
      </ResourceListElement>
    </div>
  </TanView>
</template>
