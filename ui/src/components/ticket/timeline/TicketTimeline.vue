<script setup lang="ts">
import TicketTimelineInput from '@/components/ticket/timeline/TicketTimelineInput.vue'
import TicketTimelineItem from '@/components/ticket/timeline/TicketTimelineItem.vue'
import { Card } from '@/components/ui/card'

import { Calendar as CalendarIcon } from 'lucide-vue-next'

import format from 'date-fns/format'
import { type ComputedRef, computed } from 'vue'

import type { Ticket, TimelineItem } from '@/lib/types'

const props = defineProps<{
  ticket: Ticket
  timeline?: Array<TimelineItem>
}>()

const commentsByDate: ComputedRef<Record<string, Array<TimelineItem>>> = computed(() => {
  if (!props.timeline) return {}
  const commentsByDate = props.timeline.reduce(
    (acc: Record<string, Array<TimelineItem>>, comment: TimelineItem) => {
      const date = format(new Date(comment.time), 'yyyy-MM-dd')
      if (!acc[date]) acc[date] = []
      acc[date].push(comment)
      return acc
    },
    {}
  )

  return Object.keys(commentsByDate)
    .sort()
    .reduce((acc: Record<string, Array<TimelineItem>>, date: string) => {
      acc[date] = commentsByDate[date].sort((a: TimelineItem, b: TimelineItem) => {
        return a.time > b.time ? 1 : -1
      })
      return acc
    }, {})
})
</script>

<template>
  <div class="mt-2 flex flex-col gap-2">
    <Card
      v-if="!timeline || timeline.length === 0"
      class="flex h-10 items-center p-4 text-muted-foreground"
    >
      No timeline entries added yet.
    </Card>
    <div v-else class="flex flex-col gap-2">
      <div v-for="(dateComments, date) in commentsByDate" :key="date" class="flex flex-col">
        <h2 class="flex flex-row items-center text-sm font-semibold">
          <CalendarIcon class="m-1 size-4" />
          {{ date }}
        </h2>
        <Card>
          <TicketTimelineItem
            v-for="comment in dateComments"
            :key="comment.id"
            :timelineItem="comment"
          />
        </Card>
      </div>
    </div>
    <TicketTimelineInput :ticket="ticket" class="w-full" />
  </div>
</template>
