<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import { buttonVariants } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

import { ChevronRight } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'

import { useAPI } from '@/api'
import type { ExtendedTask } from '@/client/models'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const api = useAPI()

const authStore = useAuthStore()

const {
  isPending,
  isError,
  data: tasks,
  error
} = useQuery({
  queryKey: ['tasks'],
  queryFn: (): Promise<Array<ExtendedTask>> => {
    return api
      .listTasks()
      .then((tasks) => tasks.filter((task) => task.open && task.owner === authStore.user?.id))
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
              params: { type: task.ticketType, id: task.ticket }
            }"
            :class="
              cn(
                buttonVariants({ variant: 'outline-solid', size: 'sm' }),
                'h-8 w-full sm:ml-auto sm:w-auto'
              )
            "
          >
            <span class="flex flex-row items-center text-sm text-gray-500">
              Go to {{ task.ticketName }}
              <ChevronRight class="ml-2 h-4 w-4" />
            </span>
          </RouterLink>
        </PanelListElement>
      </TanView>
    </Card>
  </div>
</template>
