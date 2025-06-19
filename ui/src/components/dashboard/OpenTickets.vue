<script setup lang="ts">
import PanelListElement from '@/components/layout/PanelListElement.vue'
import { buttonVariants } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

import { ChevronRight } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { intervalToDuration } from 'date-fns'

import { useAPI } from '@/api'
import type { ExtendedTicket } from '@/client/models'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const api = useAPI()

const authStore = useAuthStore()

const { data: tickets } = useQuery({
  queryKey: ['tickets', 'dashboard'],
  queryFn: (): Promise<Array<ExtendedTicket>> => {
    return api
      .listTickets()
      .then((tickets) =>
        tickets.filter((ticket) => ticket.open && ticket.owner === authStore.user?.id)
      )
  }
})

const age = (ticket: ExtendedTicket) => {
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
        <span class="text-muted-foreground text-sm">{{ ticket.typeSingular }}</span>
        <Separator orientation="vertical" class="hidden h-4 sm:block" />
        <span class="text-muted-foreground text-sm">Open since {{ age(ticket) }}</span>
        <RouterLink
          :to="{
            name: 'tickets',
            params: { type: ticket.type, id: ticket.id }
          }"
          :class="
            cn(
              buttonVariants({ variant: 'outline-solid', size: 'sm' }),
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
