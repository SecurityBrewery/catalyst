<script setup lang="ts">
import PanelListElement from '@/components/common/PanelListElement.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import LinkAddDialog from '@/components/ticket/link/LinkAddDialog.vue'
import LinkRemoveDialog from '@/components/ticket/link/LinkRemoveDialog.vue'

import { ref } from 'vue'

import type { Link, Ticket } from '@/lib/types'

defineProps<{
  ticket: Ticket
  links: Array<Link> | undefined
}>()

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
    <PanelListElement v-for="link in links" :key="link.id" :title="link.url" class="pr-1">
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

      <LinkRemoveDialog :ticket="ticket" :link="link" />
    </PanelListElement>
  </TicketPanel>
</template>
