<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ResourceListElement from '@/components/layout/ResourceListElement.vue'
import { Button } from '@/components/ui/button'

import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Type } from '@/client/models'

const api = useAPI()

const route = useRoute()
const router = useRouter()

const {
  isPending,
  isError,
  data: types,
  error
} = useQuery({
  queryKey: ['types'],
  queryFn: (): Promise<Array<Type>> => api.listTypes()
})

const openNew = () => {
  router.push({ name: 'types', params: { id: 'new' } })
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader title="Types">
      <div class="ml-auto">
        <Button variant="ghost" @click="openNew">New Type</Button>
      </div>
    </ColumnHeader>
    <div class="mt-2 flex flex-1 flex-col gap-2 overflow-scroll p-2 pt-0">
      <ResourceListElement
        v-for="type in types"
        :key="type.id"
        :icon="type.icon"
        :title="type.singular"
        :created="type.created"
        subtitle=""
        description=""
        :active="route.params.id === type.id"
        :to="{ name: 'types', params: { id: type.id } }"
        :open="false"
      >
        {{ type.name }}
      </ResourceListElement>
    </div>
  </TanView>
</template>
