<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ResourceListElement from '@/components/common/ResourceListElement.vue'
import PlaybookNewDialog from '@/components/playbook/PlaybookNewDialog.vue'
import { Separator } from '@/components/ui/separator'

import { useQuery } from '@tanstack/vue-query'
import { useRoute } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Playbook } from '@/lib/types'

const route = useRoute()

const {
  isPending,
  isError,
  data: playbooks,
  error
} = useQuery({
  queryKey: ['playbooks'],
  queryFn: (): Promise<Array<Playbook>> =>
    pb.collection('playbooks').getFullList({
      sort: '-created'
    })
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error" :value="playbooks">
    <div class="flex h-screen flex-col">
      <div class="flex items-center bg-background px-4 py-2">
        <h1 class="text-xl font-bold">Playbooks</h1>
        <div class="ml-auto">
          <PlaybookNewDialog />
        </div>
      </div>
      <Separator />
      <div class="mt-2 flex flex-1 flex-col gap-2 p-4 pt-0">
        <TransitionGroup name="list" appear>
          <ResourceListElement
            v-for="playbook in playbooks"
            :key="playbook.id"
            :title="playbook.name"
            :created="playbook.created"
            subtitle=""
            description=""
            :active="route.params.id === playbook.id"
            :to="{ name: 'playbooks', params: { id: playbook.id } }"
            :open="false"
          >
            {{ playbook.name }}
          </ResourceListElement>
        </TransitionGroup>
      </div>
    </div>
  </TanView>
</template>
