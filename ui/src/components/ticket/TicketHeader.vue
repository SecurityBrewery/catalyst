<script setup lang="ts">
import DynamicInput from '@/components/input/DynamicInput.vue'
import { Separator } from '@/components/ui/separator'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const name = ref(props.ticket.name)

const editNameMutation = useMutation({
  mutationFn: () =>
    pb.collection('tickets').update(props.ticket.id, {
      name: name.value
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] }),
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
