<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ResourceListElement from '@/components/layout/ResourceListElement.vue'
import { Button } from '@/components/ui/button'

import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Group } from '@/client/models'

const api = useAPI()

const route = useRoute()
const router = useRouter()

const {
  isPending,
  isError,
  data: groups,
  error
} = useQuery({
  queryKey: ['groups'],
  queryFn: (): Promise<Array<Group>> => api.listGroups()
})

const description = (group: Group): string => {
  return group.permissions.join(', ')
}

const openNew = () => {
  router.push({ name: 'groups', params: { id: 'new' } })
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader title="Groups" show-sidebar-trigger>
      <div class="ml-auto">
        <Button variant="ghost" @click="openNew">New Group</Button>
      </div>
    </ColumnHeader>
    <div class="mt-2 flex flex-1 flex-col gap-2 overflow-auto p-2 pt-0">
      <ResourceListElement
        v-for="group in groups"
        :key="group.id"
        :title="group.name"
        :created="group.created"
        subtitle=""
        :description="description(group)"
        :active="route.params.id === group.id"
        :to="{ name: 'groups', params: { id: group.id } }"
        :open="false"
      >
        {{ group.name }}
      </ResourceListElement>
    </div>
  </TanView>
</template>
