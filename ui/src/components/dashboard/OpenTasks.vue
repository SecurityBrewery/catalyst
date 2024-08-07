<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import { buttonVariants } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

import { ChevronRight } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Task } from '@/lib/types'
import { cn } from '@/lib/utils'

const {
  isPending,
  isError,
  data: tasks,
  error
} = useQuery({
  queryKey: ['tasks'],
  queryFn: (): Promise<Array<Task>> => {
    if (!pb.authStore.model) return Promise.reject('Not authenticated')
    return pb.collection('tasks').getFullList({
      sort: '-created',
      filter: pb.filter(`open = true && owner = {:owner}`, { owner: pb.authStore.model.id }),
      expand: 'owner,ticket'
    })
  }
})
</script>

<template>
  <div class="flex flex-col gap-2">
    <Card>
      <TanView :isError="isError" :isPending="isPending" :error="error">
        <div v-if="tasks && tasks.length === 0" class="p-2 text-center text-sm text-gray-500">
          No open tasks
        </div>
        <PanelListElement v-else v-for="task in tasks" :key="task.id" class="pr-1">
          <span>{{ task.name }}</span>
          <RouterLink
            :to="{
              name: 'tickets',
              params: { type: task.expand.ticket.type, id: task.expand.ticket.id }
            }"
            :class="
              cn(
                buttonVariants({ variant: 'outline', size: 'sm' }),
                'h-8 w-full sm:ml-auto sm:w-auto'
              )
            "
          >
            <span class="flex flex-row items-center text-sm text-gray-500">
              Go to {{ task.expand.ticket.name }}
              <ChevronRight class="ml-2 h-4 w-4" />
            </span>
          </RouterLink>
        </PanelListElement>
      </TanView>
    </Card>
  </div>
</template>
