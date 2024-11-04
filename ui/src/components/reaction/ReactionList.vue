<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ResourceListElement from '@/components/layout/ResourceListElement.vue'
import { Button } from '@/components/ui/button'

import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Reaction } from '@/lib/types'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()

const {
  isPending,
  isError,
  data: reactions,
  error
} = useQuery({
  queryKey: ['reactions'],
  queryFn: (): Promise<Array<Reaction>> =>
    pb.collection('reactions').getFullList({
      sort: '-created'
    })
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

onMounted(() => {
  pb.collection('reactions').subscribe('*', () => {
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
  })
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader title="Reactions">
      <div class="ml-auto">
        <Button variant="ghost" @click="openNew">New Reaction</Button>
      </div>
    </ColumnHeader>
    <div class="mt-2 flex flex-1 flex-col gap-2 p-2 pt-0">
      <TransitionGroup name="list" appear>
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
      </TransitionGroup>
    </div>
  </TanView>
</template>
