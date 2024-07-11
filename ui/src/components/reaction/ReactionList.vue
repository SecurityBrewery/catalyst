<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ResourceListElement from '@/components/common/ResourceListElement.vue'
import ReactionNewDialog from '@/components/reaction/ReactionNewDialog.vue'
import ReactionPythonNewDialog from '@/components/reaction/ReactionPythonNewDialog.vue'
import { Separator } from '@/components/ui/separator'

import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Reaction } from '@/lib/types'

const route = useRoute()

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

const subtitle = (reaction: Reaction) => {
  if (reaction.type === 'python') {
    return 'Python'
  } else if (reaction.type === 'webhook') {
    return 'Webhook'
  } else {
    return 'Unknown'
  }
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error" :value="reactions">
    <div class="flex h-screen flex-col">
      <div class="flex items-center bg-background px-4 py-2">
        <h1 class="text-xl font-bold">Reactions</h1>
        <div class="ml-auto">
          <ReactionNewDialog />
        </div>
      </div>
      <Separator />
      <div class="mt-2 flex flex-1 flex-col gap-2 p-4 pt-0">
        <TransitionGroup name="list" appear>
          <ResourceListElement
            v-for="reaction in reactions"
            :key="reaction.id"
            :title="reaction.name"
            :created="reaction.created"
            :subtitle="subtitle(reaction)"
            description=""
            :active="route.params.id === reaction.id"
            :to="{ name: 'reactions', params: { id: reaction.id, type: reaction.type } }"
            :open="false"
          >
            {{ reaction.name }}
          </ResourceListElement>
        </TransitionGroup>
      </div>
    </div>
  </TanView>
</template>
