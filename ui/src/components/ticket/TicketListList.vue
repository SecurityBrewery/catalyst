<script lang="ts" setup>
import ResourceListElement from '@/components/layout/ResourceListElement.vue'

import { useRoute } from 'vue-router'

import type { SearchTicket } from '@/lib/types'

const route = useRoute()

defineProps<{
  tickets: Array<SearchTicket>
}>()
</script>

<template>
  <div class="mt-2 flex w-full flex-1 flex-col gap-2 p-2 pt-0">
    <ResourceListElement
      v-for="item of tickets"
      :key="item.id"
      :title="item.name"
      :created="item.created"
      :subtitle="item.owner_name"
      :description="item.description ? item.description.substring(0, 300) : ''"
      :active="route.params.id === item.id"
      :to="`/tickets/${item.type}/${item.id}`"
      :open="item.open"
    />
  </div>
</template>
