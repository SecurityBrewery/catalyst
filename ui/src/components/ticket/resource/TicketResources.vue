<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import ResourceAddDialog from '@/components/ticket/resource/ResourceAddDialog.vue'
import ResourceCardContent from '@/components/ticket/resource/ResourceCardContent.vue'
import { Button } from '@/components/ui/button'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/components/ui/hover-card'

import { Trash2 } from 'lucide-vue-next'

import { type Ref, computed, ref } from 'vue'

import type { Resource, Ticket } from '@/lib/types'

const props = defineProps<{
  ticket: Ticket
  resources: Array<Resource> | undefined
}>()

const dialogOpen = ref(false)

const groups: Ref<Record<string, Array<Resource>>> = computed(() => {
  if (!props.resources) {
    return {}
  }

  const g = props.resources.reduce((acc: any, item: Resource) => {
    if (!acc[item.type]) {
      acc[item.type] = []
    }

    acc[item.type].push(item)

    return acc
  }, {})

  // sort groups by type, put 'links' last
  const keys = Object.keys(g).sort((a, b) => {
    if (a === 'other') {
      return 1
    }

    if (b === 'other') {
      return -1
    }

    return a.localeCompare(b)
  })

  return Object.fromEntries(keys.map((key) => [key, g[key]]))
})
</script>

<template>
  <TicketPanel title="Resources" @add="dialogOpen = true">
    <ResourceAddDialog v-model="dialogOpen" :ticket="ticket" />
    <template v-for="(items, name) in groups" :key="name">
      <h3 class="p-2 text-sm font-semibold bg-border">
        {{ name }}
      </h3>
      <div
        v-if="!resources || resources.length === 0"
        class="flex h-10 items-center p-4 text-muted-foreground"
      >
        No links added yet.
      </div>
      <HoverCard v-for="resource in items" :key="resource.id" :openDelay="200">
        <HoverCardTrigger as-child>
          <PanelListElement :title="resource.url" class="flex-row items-center pr-1">
            <span v-if="resource.url" class="flex-1 overflow-hidden overflow-ellipsis">
              <a
                :href="resource.url"
                target="_blank"
                class="flex flex-1 items-center overflow-hidden"
              >
                <span class="mr-2 text-blue-500 underline">
                  {{ resource.name }}
                </span>

                <div
                  class="flex-1 overflow-hidden overflow-ellipsis text-nowrap text-sm text-muted-foreground"
                >
                  {{ resource.url }}
                </div>
              </a>
            </span>
            <span v-else class="flex-1 overflow-hidden overflow-ellipsis">
              <span class="mr-2">
                {{ resource.name }}
              </span>
            </span>

            <DeleteDialog
              v-if="resource"
              collection="resources"
              :id="resource.id"
              :name="resource.name"
              singular="Resource"
              :queryKey="['tickets', ticket.id]"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
        </HoverCardTrigger>
        <HoverCardContent>
          <ResourceCardContent :resource="resource" />
        </HoverCardContent>
      </HoverCard>
    </template>
  </TicketPanel>
</template>
