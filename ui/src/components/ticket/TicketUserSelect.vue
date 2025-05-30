<script setup lang="ts">
import UserSelect from '@/components/common/UserSelect.vue'

import { useMutation, useQueryClient } from '@tanstack/vue-query'

import { api } from '@/api'
import type { ExtendedTicket, Ticket, User } from '@/client/models'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: ExtendedTicket
}>()

const setTicketOwnerMutation = useMutation({
  mutationFn: (userID: string): Promise<Ticket> =>
    api.updateTicket({ id: props.ticket.id, ticketUpdate: { owner: userID } }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets'] }),
  onError: handleError
})

const update = (user: User) => setTicketOwnerMutation.mutate(user.id)
</script>

<template>
  <UserSelect :userID="ticket.owner" :userName="ticket.ownerName" @select="update" />
</template>
