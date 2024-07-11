<script setup lang="ts">
import Icon from '@/components/Icon.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import TicketCloseDialog from '@/components/ticket/TicketCloseDialog.vue'
import TicketDeleteDialog from '@/components/ticket/TicketDeleteDialog.vue'
import TicketUserSelect from '@/components/ticket/TicketUserSelect.vue'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { toast } from '@/components/ui/toast'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import { Check, CircleDot, Repeat } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket, Type } from '@/lib/types'

const queryClient = useQueryClient()
const router = useRouter()

const props = defineProps<{
  ticket: Ticket
}>()

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

const changeTypeMutation = useMutation({
  mutationFn: (typeID: string): Promise<Ticket> =>
    pb.collection('tickets').update(props.ticket.id, {
      type: typeID
    }),
  onSuccess: (data: Ticket) => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    router.push({ name: 'tickets', params: { type: data.type, id: props.ticket.id } })
  },
  onError: (error) => {
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
  }
})

const closeTicketMutation = useMutation({
  mutationFn: (): Promise<Ticket> =>
    pb.collection('tickets').update(props.ticket.id, {
      open: !props.ticket.open
    }),
  onSuccess: (data: Ticket) => {
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    if (!data.open) {
      router.push({ name: 'tickets', params: { type: props.ticket.expand.type.id } })
    }
  },
  onError: (error) => {
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
  }
})

const otherTypes = computed(() => types.value?.filter((t) => t.id !== props.ticket.expand.type.id))

const closeTicketDialogOpen = ref(false)
</script>

<template>
  <div class="flex items-center justify-between bg-background p-2">
    <div class="flex items-center gap-2">
      <Tooltip>
        <TooltipTrigger as-child>
          <div>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button variant="outline" :disabled="!ticket">
                  <Icon :name="ticket.expand.type.icon" class="mr-2 size-4" />
                  {{ ticket.expand.type.singular }}
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
            <TicketUserSelect :key="ticket.owner" :uID="ticket.owner" :ticket="ticket" />
          </div>
        </TooltipTrigger>
        <TooltipContent>Change User</TooltipContent>
      </Tooltip>
    </div>
    <DeleteDialog
      v-if="ticket"
      :collection="'tickets'"
      :id="ticket.id"
      :name="ticket.name"
      singular="Ticket"
      :to="{ name: 'tickets' }"
      :queryKey="['tickets']"
    />
  </div>
</template>
