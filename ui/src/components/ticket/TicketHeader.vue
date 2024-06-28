<script setup lang="ts">
import DynamicInput from '@/components/input/DynamicInput.vue'
import { Separator } from '@/components/ui/separator'
import { toast } from '@/components/ui/toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Ticket, Type } from '@/lib/types'

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
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
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

  <div class="flex flex-row space-x-2 px-1 text-xs">
    <div class="flex items-center gap-1 text-muted-foreground">
      Created:
      {{ format(new Date(ticket.created), 'PPpp') }}
    </div>
    <Separator orientation="vertical" />
    <div class="flex items-center gap-1 text-muted-foreground">
      Updated:
      {{ format(new Date(ticket.updated), 'PPpp') }}
    </div>
  </div>
</template>
