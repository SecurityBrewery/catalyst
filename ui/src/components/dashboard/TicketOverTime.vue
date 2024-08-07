<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import { LineChart } from '@/components/ui/chart-line'

import { useQuery } from '@tanstack/vue-query'
import { getWeek } from 'date-fns'
import { computed } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Task } from '@/lib/types'

const {
  isPending,
  isError,
  data: tickets,
  error
} = useQuery({
  queryKey: ['tickets'],
  queryFn: (): Promise<Array<Task>> =>
    pb.collection('tickets').getFullList({
      sort: '-created',
      expand: 'owner,type'
    })
})

const ticketsPerWeek = computed(() => {
  if (!tickets.value) return []

  const weeks = tickets.value.reduce(
    (acc, ticket) => {
      const week = getWeek(ticket.created)
      acc[week] = (acc[week] || 0) + 1
      return acc
    },
    {} as Record<number, number>
  )

  return Object.entries(weeks).map(([week, count]) => ({
    week: parseInt(week),
    count
  }))
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <LineChart class="h-40" :data="ticketsPerWeek" index="week" :categories="['count']" />
  </TanView>
</template>
