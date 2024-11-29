<script setup lang="ts">
import PanelListElement from '@/components/layout/PanelListElement.vue'
import { buttonVariants } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

import { ChevronRight } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { intervalToDuration } from 'date-fns'

import { pb } from '@/lib/pocketbase'
import type { Ticket } from '@/lib/types'
import { cn } from '@/lib/utils'

const {
  isPending,
  isError,
  data: tickets,
  error
} = useQuery({
  queryKey: ['tickets', 'dashboard'],
  queryFn: (): Promise<Array<Ticket>> => {
    if (!pb.authStore.model) return Promise.reject('Not authenticated')
    return pb.collection('tickets').getFullList({
      sort: '-created',
      filter: pb.filter(`open = true && owner = {:owner}`, { owner: pb.authStore.model.id }),
      expand: 'owner,type'
    })
  }
})

const age = (ticket: Ticket) => {
  const days = intervalToDuration({ start: new Date(ticket.created), end: new Date() }).days

  if (!days) return 'today'
  if (days === 1) return 'yesterday'

  return `${days} days`
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <Card>
      <div v-if="tickets && tickets.length === 0" class="p-2 text-center text-sm text-gray-500">
        No open tickets
      </div>
      <PanelListElement v-else v-for="ticket in tickets" :key="ticket.id" class="gap-2 pr-1">
        <span>{{ ticket.name }}</span>
        <Separator orientation="vertical" class="hidden h-4 sm:block" />
        <span class="text-sm text-muted-foreground">{{ ticket.expand.type.singular }}</span>
        <Separator orientation="vertical" class="hidden h-4 sm:block" />
        <span class="text-sm text-muted-foreground">Open since {{ age(ticket) }}</span>
        <RouterLink
          :to="{
            name: 'tickets',
            params: { type: ticket.type, id: ticket.id }
          }"
          :class="
            cn(
              buttonVariants({ variant: 'outline', size: 'sm' }),
              'h-8 w-full sm:ml-auto sm:w-auto'
            )
          "
        >
          <span class="flex flex-row items-center text-sm text-gray-500">
            Go to {{ ticket.name }}
            <ChevronRight class="ml-2 h-4 w-4" />
          </span>
        </RouterLink>
      </PanelListElement>
    </Card>
  </div>
</template>
