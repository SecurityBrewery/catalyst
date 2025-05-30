<script setup lang="ts">
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import { Check, Repeat } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { api } from '@/api'
import type { ExtendedTicket, Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()
const router = useRouter()

const props = defineProps<{
  ticket: ExtendedTicket
}>()

const resolution = ref(props.ticket.resolution)

const closeTicketMutation = useMutation({
  mutationFn: (): Promise<Ticket> =>
    api.updateTicket({
      id: props.ticket.id,
      ticketUpdate: { open: !props.ticket.open, resolution: resolution.value }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    router.push({ name: 'tickets', params: { type: props.ticket.type } })
  },
  onError: handleError
})

const closeButtonDisabled = false // computed(() => !props.ticket.open || message.value == '')
</script>

<template>
  <ColumnHeader nowrap hideSeparator>
    <Input v-if="ticket.open" v-model="resolution" placeholder="Closing reason" />
    <div v-else class="flex-1">
      <p class="ml-2 text-gray-500">Closed: {{ ticket.resolution }}</p>
    </div>
    <Button
      @click="closeTicketMutation.mutate"
      :disabled="closeButtonDisabled"
      :variant="closeButtonDisabled ? 'secondary' : 'default'"
    >
      <Check v-if="ticket.open" class="mr-2 h-4 w-4" />
      <Repeat v-else class="mr-2 h-4 w-4" />
      {{
        ticket?.open ? 'Close ' + props.ticket.typeSingular : 'Reopen ' + props.ticket.typeSingular
      }}
    </Button>
  </ColumnHeader>
</template>
