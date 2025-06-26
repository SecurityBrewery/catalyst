<script setup lang="ts">
import DynamicInput from '@/components/input/DynamicInput.vue'
import { Separator } from '@/components/ui/separator'
import { useToast } from '@/components/ui/toast/use-toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

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
  onSuccess: (data: Ticket) => {
    toast({
      title: 'Ticket updated',
      description: 'The ticket name has been updated'
    })
    queryClient.invalidateQueries({ queryKey: ['tickets', data.id] })
  },
  onError: handleError('Failed to update ticket name')
})

const updateName = (value: string) => {
  name.value = value
  editNameMutation.mutate()
}
</script>

<template>
  <span class="text-4xl font-bold">
    <DynamicInput
      id="name"
      :modelValue="ticket.name"
      @update:modelValue="updateName"
      class="-mx-1"
    />
  </span>

  <div class="text-muted-foreground flex flex-col items-stretch gap-1 text-xs md:h-4 md:flex-row">
    <div>
      Created:
      {{ format(ticket.created, 'PPpp') }}
    </div>
    <Separator orientation="vertical" class="hidden md:block" />
    <div>
      Updated:
      {{ format(ticket.updated, 'PPpp') }}
    </div>
  </div>
</template>
