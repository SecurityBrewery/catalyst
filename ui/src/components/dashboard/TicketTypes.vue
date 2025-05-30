<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import { DonutChart } from '@/components/ui/chart-donut'

import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

import { api } from '@/api'
import type { Sidebar } from '@/client/models'

const {
  isPending,
  isError,
  data: types,
  error
} = useQuery({
  queryKey: ['sidebar'],
  queryFn: (): Promise<Array<Sidebar>> => api.getSidebar()
})

const namedTypes = computed(() => {
  if (!types.value) return []
  return types.value.map((type) => {
    return {
      plural: type.plural,
      name: type.plural, // fixes the donut chart, which always expects "name" as the index field
      count: type.count
    }
  })
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <div v-if="namedTypes" class="flex flex-1 items-center">
      <DonutChart index="plural" type="donut" category="count" :data="namedTypes" />
    </div>
  </TanView>
</template>
