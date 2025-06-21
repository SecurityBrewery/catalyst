<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import LinkAddDialog from '@/components/ticket/link/LinkAddDialog.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { Link, Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  ticket: Ticket
  links: Array<Link> | undefined
}>()

const deleteMutation = useMutation({
  mutationFn: (id: string) => api.deleteLink({ id }),
  onSuccess: (data, id) => {
    queryClient.removeQueries({ queryKey: ['links', id] })
    queryClient.invalidateQueries({ queryKey: ['links', props.ticket.id] })
    toast({
      title: 'Link deleted',
      description: 'The link has been deleted successfully'
    })
  },
  onError: handleError('Failed to delete link')
})

const dialogOpen = ref(false)
</script>

<template>
  <TicketPanel title="Links" @add="dialogOpen = true">
    <LinkAddDialog v-model="dialogOpen" :ticket="ticket" />
    <div
      v-if="!links || links.length === 0"
      class="text-muted-foreground flex h-10 items-center p-4"
    >
      No links added yet.
    </div>
    <PanelListElement
      v-for="link in links"
      :key="link.id"
      :title="link.url"
      class="flex-row items-center pr-1"
    >
      <a :href="link.url" target="_blank" class="flex flex-1 items-center overflow-hidden">
        <span class="mr-2 text-blue-500 underline">
          {{ link.name }}
        </span>

        <div class="text-muted-foreground flex-1 overflow-hidden text-sm text-nowrap text-ellipsis">
          {{ link.url }}
        </div>
      </a>

      <DeleteDialog
        v-if="link"
        :name="link.name"
        singular="Link"
        @delete="deleteMutation.mutate(link.id)"
      >
        <Button variant="ghost" size="icon" class="h-8 w-8">
          <Trash2 class="size-4" />
          <span class="sr-only">Delete Link</span>
        </Button>
      </DeleteDialog>
    </PanelListElement>
  </TicketPanel>
</template>
