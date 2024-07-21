<script lang="ts" setup>
import ResourceListElement from '@/components/common/ResourceListElement.vue'

import { useRoute } from 'vue-router'

import type { Ticket } from '@/lib/types'

const route = useRoute()

defineProps<{
  tickets: Array<Ticket>
}>()
</script>

<template>
  <div class="mt-2 flex w-full flex-1 flex-col gap-2 p-4 pt-0">
    <ResourceListElement
      v-for="item of tickets"
      :key="item.id"
      :title="item.name"
      :created="item.created"
      :subtitle="item.expand.owner ? item.expand.owner.name : ''"
      :description="item.description ? item.description.substring(0, 300) : ''"
      :active="route.params.id === item.id"
      :to="`/tickets/${item.expand.type.id}/${item.id}`"
      :open="item.open"
    />
  </div>
</template>
