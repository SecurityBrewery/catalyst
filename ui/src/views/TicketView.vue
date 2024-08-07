<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import TicketDisplay from '@/components/ticket/TicketDisplay.vue'
import TicketList from '@/components/ticket/TicketList.vue'

import { useQuery } from '@tanstack/vue-query'
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Type } from '@/lib/types'

const route = useRoute()
const router = useRouter()

const id = computed(() => route.params.id as string)
const type = computed(() => route.params.type as string)

const {
  isPending,
  isError,
  data: selectedType,
  error,
  refetch
} = useQuery({
  queryKey: ['types', type.value],
  queryFn: (): Promise<Type> => pb.collection('types').getOne(type.value)
})

watch(
  () => type.value,
  () => refetch()
)

onMounted(() => {
  if (!pb.authStore.model) {
    router.push({ name: 'login' })
  }
})
</script>

<template>
  <ThreeColumn :show-details="!!id">
    <template #list>
      <TanView :isError="isError" :isPending="isPending" :error="error">
        <TicketList v-if="selectedType" :key="selectedType.id" :selectedType="selectedType" />
      </TanView>
    </template>
    <template #single>
      <TanView :isError="isError" :isPending="isPending" :error="error">
        <ColumnBody v-if="!id" class="items-center justify-center text-lg text-gray-500">
          No ticket selected
        </ColumnBody>
        <TicketDisplay v-else-if="selectedType" :key="id" :selectedType="selectedType" />
      </TanView>
    </template>
  </ThreeColumn>
</template>
