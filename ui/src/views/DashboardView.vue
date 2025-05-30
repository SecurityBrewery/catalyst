<script setup lang="ts">
import OpenTasks from '@/components/dashboard/OpenTasks.vue'
import OpenTickets from '@/components/dashboard/OpenTickets.vue'
import TicketOverTime from '@/components/dashboard/TicketOverTime.vue'
import TicketTypes from '@/components/dashboard/TicketTypes.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import TwoColumn from '@/components/layout/TwoColumn.vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { ExternalLink } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

import { api } from '@/api'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()

const {
  isPending,
  isError,
  data: dashboardCounts,
  error
} = useQuery({
  queryKey: ['dashboard_counts'],
  queryFn: (): Promise<Array<any>> => {
    return api.getDashboardCounts()
  }
})

const count = (id: string) => {
  if (!dashboardCounts.value) return 0

  const s = dashboardCounts.value.filter((stat) => stat.id === id)
  if (s.length === 0) return 0

  return s[0].count
}

onMounted(() => {
  if (!authStore.user) {
    router.push({ name: 'login' })
  }
})
</script>

<template>
  <TwoColumn>
    <ColumnHeader title="Dashboard" />
    <ColumnBody>
      <ColumnBodyContainer
        class="grid grid-cols-1 grid-rows-[100px_100px_100px_100px] md:grid-cols-2 md:grid-rows-[100px_100px] xl:grid-cols-4 xl:grid-rows-[100px]"
      >
        <Card>
          <CardHeader>
            <CardTitle>{{ count('tasks') }}</CardTitle>
            <CardDescription>Tasks</CardDescription>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>{{ count('tickets') }}</CardTitle>
            <CardDescription>Tickets</CardDescription>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>{{ count('users') }}</CardTitle>
            <CardDescription>Users</CardDescription>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>{{ count('reactions') }}</CardTitle>
            <CardDescription>Reactions</CardDescription>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle> Catalyst</CardTitle>
          </CardHeader>
          <CardContent class="flex flex-1 flex-col gap-1">
            <a
              href="https://catalyst.security-brewery.com/docs/category/catalyst-handbook"
              target="_blank"
              class="flex items-center rounded border p-2 text-blue-500 hover:bg-accent"
            >
              Open Catalyst Handbook
              <ExternalLink class="ml-2 h-4 w-4" />
            </a>
            <a
              href="/_/"
              target="_blank"
              class="flex items-center rounded border p-2 text-blue-500 hover:bg-accent"
            >
              Open Admin Interface
              <ExternalLink class="ml-2 h-4 w-4" />
            </a>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle> Tickets by Type</CardTitle>
          </CardHeader>
          <CardContent>
            <TicketTypes />
          </CardContent>
        </Card>
        <Card class="xl:col-span-2">
          <CardHeader>
            <CardTitle>Tickets Per Week</CardTitle>
          </CardHeader>
          <CardContent>
            <TicketOverTime />
          </CardContent>
        </Card>
        <Card class="xl:col-span-2">
          <CardHeader>
            <CardTitle>Your Open Tickets</CardTitle>
          </CardHeader>
          <CardContent>
            <OpenTickets />
          </CardContent>
        </Card>
        <Card class="xl:col-span-2">
          <CardHeader>
            <CardTitle>Your Open Tasks</CardTitle>
          </CardHeader>
          <CardContent>
            <OpenTasks />
          </CardContent>
        </Card>
      </ColumnBodyContainer>
    </ColumnBody>
  </TwoColumn>
</template>
