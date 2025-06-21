<script setup lang="ts">
import Icon from '@/components/Icon.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import TicketCloseDialog from '@/components/ticket/TicketCloseDialog.vue'
import TicketUserSelect from '@/components/ticket/TicketUserSelect.vue'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { useToast } from '@/components/ui/toast/use-toast'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import { Check, ChevronLeft, CircleDot, Repeat } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { ExtendedTicket, Ticket } from '@/client/models'
import type { Type } from '@/client/models/Type'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()
const { toast } = useToast()

const props = defineProps<{
  ticket: ExtendedTicket
}>()

const { data: types } = useQuery({
  queryKey: ['types'],
  queryFn: (): Promise<Array<Type>> => api.listTypes()
})

const changeTypeMutation = useMutation({
  mutationFn: (typeID: string): Promise<Ticket> =>
    api.updateTicket({ id: props.ticket.id, ticketUpdate: { type: typeID } }),
  onSuccess: (data: Ticket) => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    queryClient.invalidateQueries({ queryKey: ['sidebar'] })
    toast({
      title: 'Type changed',
      description: 'The ticket type has been updated'
    })
    router.push({ name: 'tickets', params: { type: data.type, id: props.ticket.id } })
  },
  onError: handleError('Failed to change ticket type')
})

const closeTicketMutation = useMutation({
  mutationFn: (): Promise<Ticket> =>
    api.updateTicket({ id: props.ticket.id, ticketUpdate: { open: !props.ticket.open } }),
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

const ticketType = computed(() => types.value?.find((t) => t.id === props.ticket.type))

const otherTypes = computed(() => types.value?.filter((t) => t.id !== props.ticket.type))

const closeTicketDialogOpen = ref(false)

const deleteMutation = useMutation({
  mutationFn: () => api.deleteTicket({ id: props.ticket.id }),
  onSuccess: () => {
    queryClient.removeQueries({ queryKey: ['tickets', props.ticket.id] })
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    toast({
      title: 'Ticket deleted',
      description: 'The ticket has been deleted successfully'
    })
    router.push({ name: 'tickets' })
  },
  onError: handleError('Failed to delete ticket')
})
</script>

<template>
  <ColumnHeader>
    <Button
      @click="router.push({ name: 'tickets', params: { type: ticket.type } })"
      variant="outline"
      class="sm:hidden"
    >
      <ChevronLeft class="mr-2 size-4" />
      Back
    </Button>
    <Tooltip>
      <TooltipTrigger as-child>
        <div>
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <Button variant="outline" :disabled="!ticket">
                <Icon :name="ticketType?.icon ?? ''" class="mr-2 size-4" />
                {{ ticketType?.singular }}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem
                v-for="type in otherTypes"
                :key="type.id"
                class="cursor-pointer"
                @click="changeTypeMutation.mutate(type.id)"
              >
                <Icon :name="type.icon" class="mr-2 size-4" />
                Convert to {{ type.singular }}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </TooltipTrigger>
      <TooltipContent>Change Type</TooltipContent>
    </Tooltip>
    <TicketCloseDialog v-model="closeTicketDialogOpen" :ticket="ticket" />
    <Tooltip>
      <TooltipTrigger as-child>
        <div>
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <Button variant="outline" :disabled="!ticket">
                <CircleDot v-if="ticket.open" class="mr-2 h-4 w-4" />
                <Check v-else class="mr-2 h-4 w-4" />
                {{ ticket?.open ? 'Open' : 'Closed' }}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem
                v-if="ticket.open"
                class="cursor-pointer"
                @click="closeTicketDialogOpen = true"
              >
                <Check class="mr-2 size-4" />
                Close Ticket
              </DropdownMenuItem>
              <DropdownMenuItem v-else class="cursor-pointer" @click="closeTicketMutation.mutate">
                <Repeat class="mr-2 size-4" />
                Reopen Ticket
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </TooltipTrigger>
      <TooltipContent>Change Status</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <div>
          <TicketUserSelect :key="ticket.owner" :ticket="ticket" />
        </div>
      </TooltipTrigger>
      <TooltipContent>Change User</TooltipContent>
    </Tooltip>
    <div class="-mx-1 flex-1" />
    <DeleteDialog
      v-if="ticket"
      :name="ticket.name"
      singular="Ticket"
      @delete="deleteMutation.mutate"
    />
  </ColumnHeader>
</template>
