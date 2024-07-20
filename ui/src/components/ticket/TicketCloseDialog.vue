<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Textarea } from '@/components/ui/textarea'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()
const router = useRouter()

const model = defineModel<boolean>()

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
  onSuccess: (data: Ticket) => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    if (!data.open) {
      router.push({ name: 'tickets', params: { type: props.ticket.expand.type.id } })
    }
  },
  onError: handleError
})
</script>

<template>
  <Dialog v-model:open="model">
    <DialogContent>
      <DialogHeader>
        <DialogTitle> Close Ticket "{{ props.ticket.name }}"</DialogTitle>
        <DialogDescription> Are you sure you want to close this ticket?</DialogDescription>
      </DialogHeader>

      <Textarea v-model="resolution" placeholder="Closing reason" />

      <DialogFooter class="mt-2 sm:justify-start">
        <Button type="button" variant="default" @click="closeTicketMutation.mutate()">
          Close
        </Button>
        <DialogClose as-child>
          <Button type="button" variant="secondary">Cancel</Button>
        </DialogClose>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
