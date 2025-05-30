<script setup lang="ts">
import DynamicInput from '@/components/input/DynamicInput.vue'
import { Separator } from '@/components/ui/separator'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { api } from '@/api'
import type { Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const name = ref(props.ticket.name)

const editNameMutation = useMutation({
  mutationFn: (): Promise<Ticket> =>
    api.updateTicket({
      id: props.ticket.id,
      ticketUpdate: {
        name: name.value
      }
    }),
  onSuccess: (data: Ticket) => queryClient.invalidateQueries({ queryKey: ['tickets', data.id] }),
  onError: handleError
})

const updateName = (value: string) => {
  name.value = value
  editNameMutation.mutate()
}
</script>

<template>
  <span class="text-4xl font-bold">
    <DynamicInput :modelValue="ticket.name" @update:modelValue="updateName" class="-mx-1" />
  </span>

  <div class="flex flex-col items-stretch gap-1 text-xs text-muted-foreground md:h-4 md:flex-row">
    <div>
      Created:
      {{ format(new Date(ticket.created), 'PPpp') }}
    </div>
    <Separator orientation="vertical" class="hidden md:block" />
    <div>
      Updated:
      {{ format(new Date(ticket.updated), 'PPpp') }}
    </div>
  </div>
</template>
