<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { toast } from '@/components/ui/toast'

import { MoreVertical, Trash2 } from 'lucide-vue-next'

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

const isOpen = ref(false)

const deleteTicketMutation = useMutation({
  mutationFn: (): Promise<boolean> => pb.collection('tickets').delete(props.ticket.id),
  onSuccess: () => {
    router.push({ name: 'tickets', params: { type: props.ticket.expand.type.id } })
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    isOpen.value = false
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" :disabled="!ticket">
          <MoreVertical class="size-4" />
          <span class="sr-only">More</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem class="cursor-pointer" as-child>
          <DialogTrigger class="w-full">
            <Trash2 class="mr-2 h-4 w-4" />
            Delete Ticket
          </DialogTrigger>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <DialogContent>
      <DialogHeader>
        <DialogTitle> Delete Ticket "{{ props.ticket.name }}"</DialogTitle>
        <DialogDescription> Are you sure you want to delete this ticket?</DialogDescription>
      </DialogHeader>

      <DialogFooter class="mt-2">
        <DialogClose as-child>
          <Button type="button" variant="secondary"> Cancel</Button>
        </DialogClose>
        <Button type="button" variant="destructive" @click="deleteTicketMutation.mutate()">
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
