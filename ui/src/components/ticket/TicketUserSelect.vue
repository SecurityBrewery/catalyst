<script setup lang="ts">
import UserSelect from '@/components/common/UserSelect.vue'
import { useToast } from '@/components/ui/toast/use-toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'

import { useAPI } from '@/api'
import type { ExtendedTicket, Ticket, User } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  ticket: ExtendedTicket
}>()

const setTicketOwnerMutation = useMutation({
  mutationFn: (userID: string): Promise<Ticket> =>
    api.updateTicket({ id: props.ticket.id, ticketUpdate: { owner: userID } }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    toast({
      title: 'Owner changed',
      description: 'The ticket owner has been updated'
    })
  },
  onError: handleError('Failed to update ticket owner')
})

const update = (user: User) => setTicketOwnerMutation.mutate(user.id)
</script>

<template>
  <UserSelect :userID="ticket.owner" :userName="ticket.ownerName" @select="update" />
</template>
