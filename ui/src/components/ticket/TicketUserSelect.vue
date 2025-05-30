<script setup lang="ts">
import UserSelect from '@/components/common/UserSelect.vue'
import { Button } from '@/components/ui/button'

import { User2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'

import { api } from '@/api'
import type { ExtendedTicket, Ticket } from '@/client/models'
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

const update = (userID: string) => setTicketOwnerMutation.mutate(userID)
</script>

<template>
  <UserSelect v-if="!ticket.owner" @update:modelValue="update">
    <Button variant="outline" role="combobox">
      <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
      Unassigned
    </Button>
  </UserSelect>
  <UserSelect v-else :modelValue="ticket.owner" @update:modelValue="update">
    <Button variant="outline" role="combobox">
      <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
      {{ ticket.ownerName }}
    </Button>
  </UserSelect>
</template>
