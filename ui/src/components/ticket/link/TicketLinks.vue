<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import LinkAddDialog from '@/components/ticket/link/LinkAddDialog.vue'
import { Button } from '@/components/ui/button'

import { Trash2 } from 'lucide-vue-next'

import { useMutation } from '@tanstack/vue-query'
import { useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { api } from '@/api'
import type { Link, Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  links: Array<Link> | undefined
}>()

const deleteMutation = useMutation({
  mutationFn: () => {
    return api.deleteLink({ id: props.ticket.id })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['links', props.ticket.id] })
  },
  onError: handleError
})

const dialogOpen = ref(false)
</script>

<template>
  <TicketPanel title="Links" @add="dialogOpen = true">
    <LinkAddDialog v-model="dialogOpen" :ticket="ticket" />
    <div
      v-if="!links || links.length === 0"
      class="flex h-10 items-center p-4 text-muted-foreground"
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

        <div
          class="flex-1 overflow-hidden overflow-ellipsis text-nowrap text-sm text-muted-foreground"
        >
          {{ link.url }}
        </div>
      </a>

      <DeleteDialog v-if="link" :name="link.name" singular="Link" @delete="deleteMutation.mutate">
        <Button variant="ghost" size="icon" class="h-8 w-8">
          <Trash2 class="size-4" />
        </Button>
      </DeleteDialog>
    </PanelListElement>
  </TicketPanel>
</template>
