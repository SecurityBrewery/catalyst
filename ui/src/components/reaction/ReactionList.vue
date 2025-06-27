<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ResourceListElement from '@/components/layout/ResourceListElement.vue'
import { Button } from '@/components/ui/button'

import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Reaction } from '@/client/models'

const api = useAPI()

const route = useRoute()
const router = useRouter()

const {
  isPending,
  isError,
  data: reactions,
  error
} = useQuery({
  queryKey: ['reactions'],
  queryFn: (): Promise<Array<Reaction>> => api.listReactions()
})

const subtitle = (reaction: Reaction) =>
  triggerNiceName(reaction) + ' to ' + reactionNiceName(reaction)

const triggerNiceName = (reaction: Reaction) => {
  if (reaction.trigger === 'schedule') {
    return 'Schedule'
  } else if (reaction.trigger === 'hook') {
    return 'Collection Hook'
  } else if (reaction.trigger === 'webhook') {
    return 'HTTP / Webhook'
  } else {
    return 'Unknown'
  }
}

const reactionNiceName = (reaction: Reaction) => {
  if (reaction.action === 'python') {
    return 'Python'
  } else if (reaction.action === 'webhook') {
    return 'HTTP / Webhook'
  } else {
    return 'Unknown'
  }
}

const openNew = () => {
  router.push({ name: 'reactions', params: { id: 'new' } })
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader title="Reactions" show-sidebar-trigger>
      <div class="ml-auto">
        <Button variant="ghost" @click="openNew">New Reaction</Button>
      </div>
    </ColumnHeader>
    <div class="mt-2 flex flex-1 flex-col gap-2 overflow-auto p-2 pt-0">
      <ResourceListElement
        v-for="reaction in reactions"
        :key="reaction.id"
        :title="reaction.name"
        :created="reaction.created"
        :subtitle="subtitle(reaction)"
        description=""
        :active="route.params.id === reaction.id"
        :to="{ name: 'reactions', params: { id: reaction.id } }"
        :open="false"
      >
        {{ reaction.name }}
      </ResourceListElement>
    </div>
  </TanView>
</template>
