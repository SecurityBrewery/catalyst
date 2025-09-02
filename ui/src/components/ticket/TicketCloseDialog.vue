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
import { useToast } from '@/components/ui/toast/use-toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()
const { toast } = useToast()

const model = defineModel<boolean>()

const props = defineProps<{
  ticket: Ticket
}>()

const resolution = ref(props.ticket.resolution)

const closeTicketMutation = useMutation({
  mutationFn: (): Promise<Ticket> =>
    api.updateTicket({
      id: props.ticket.id,
      ticketUpdate: {
        open: !props.ticket.open,
        resolution: resolution.value
      }
    }),
  onSuccess: (data: Ticket) => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    toast({
      title: data.open ? 'Ticket reopened' : 'Ticket closed',
      description: data.open
        ? 'The ticket has been reopened successfully'
        : 'The ticket has been closed successfully'
    })
    if (!data.open) {
      router.push({ name: 'tickets', params: { type: props.ticket.type } })
    }
  },
  onError: handleError('Failed to update ticket status')
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
