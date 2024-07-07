<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { toast } from '@/components/ui/toast'

import { Check, Repeat } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket } from '@/lib/types'

const queryClient = useQueryClient()
const router = useRouter()

const props = defineProps<{
  ticket: Ticket
}>()

const resolution = ref(props.ticket.resolution)

const closeTicketMutation = useMutation({
  mutationFn: (): Promise<Ticket> =>
    pb.collection('tickets').update(props.ticket.id, {
      open: !props.ticket.open,
      resolution: resolution.value
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    router.push({ name: 'tickets', params: { type: props.ticket.expand.type.id } })
  },
  onError: (error) => {
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
  }
})

const closeButtonDisabled = false // computed(() => !props.ticket.open || message.value == '')
</script>

<template>
  <div class="flex items-center justify-between gap-2 bg-background p-2">
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
        ticket?.open
          ? 'Close ' + props.ticket.expand.type.singular
          : 'Reopen ' + props.ticket.expand.type.singular
      }}
    </Button>
  </div>
</template>
