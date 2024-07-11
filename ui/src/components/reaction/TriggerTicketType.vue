<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'

import { useQuery } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Type } from '@/lib/types'

interface TriggerTicketTypeData {
  tickettype: Array<string>
}

const props = defineProps<{
  modelValue: TriggerTicketTypeData
}>()

const emit = defineEmits(['update:modelValue'])

const updateTicketType = (id: string, value: boolean) => {
  let tickettype = props.modelValue.tickettype || []

  if (value) {
    tickettype = [...tickettype, id]
  } else {
    tickettype = tickettype.filter((t) => t !== id)
  }

  emit('update:modelValue', { ...props.modelValue, tickettype })
}

const {
  isPending,
  isError,
  data: types,
  error
} = useQuery({
  queryKey: ['types'],
  queryFn: (): Promise<Array<Type>> =>
    pb.collection('types').getFullList({
      sort: '-created'
    })
})
</script>

<template>
  <p class="py-4 text-sm text-muted-foreground">Add controls to tickets of the specified types.</p>

  <div class="mt-2 flex items-center space-x-2" v-for="type in types" :key="type.id">
    <Switch
      :id="type.id"
      :checked="modelValue.tickettype && modelValue.tickettype.includes(type.id)"
      @update:checked="(value) => updateTicketType(type.id, value)"
    />
    <Label :for="type.id">Show controls for {{ type.plural }}</Label>
  </div>
</template>
