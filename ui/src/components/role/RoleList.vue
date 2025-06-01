<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ResourceListElement from '@/components/layout/ResourceListElement.vue'
import { Button } from '@/components/ui/button'

import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Role } from '@/client/models'

const api = useAPI()

const route = useRoute()
const router = useRouter()

const {
  isPending,
  isError,
  data: roles,
  error
} = useQuery({
  queryKey: ['roles'],
  queryFn: (): Promise<Array<Role>> => api.listRoles()
})

const description = (role: Role): string => {
  return role.permissions.join(', ')
}

const openNew = () => {
  router.push({ name: 'roles', params: { id: 'new' } })
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader title="Roles">
      <div class="ml-auto">
        <Button variant="ghost" @click="openNew">New Role</Button>
      </div>
    </ColumnHeader>
    <div class="mt-2 flex flex-1 flex-col gap-2 overflow-scroll p-2 pt-0">
      <ResourceListElement
        v-for="role in roles"
        :key="role.id"
        :title="role.name"
        :created="role.created"
        subtitle=""
        :description="description(role)"
        :active="route.params.id === role.id"
        :to="{ name: 'roles', params: { id: role.id } }"
        :open="false"
      >
        {{ role.name }}
      </ResourceListElement>
    </div>
  </TanView>
</template>
