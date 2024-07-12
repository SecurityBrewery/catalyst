<script setup lang="ts">
import MDEditor from '@/components/input/MDEditor.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Calendar } from '@/components/ui/v-calendar'

import { Calendar as CalendarIcon, Plus } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Ticket, TimelineItem } from '@/lib/types'
import { cn, handleError } from '@/lib/utils'

const props = defineProps<{
  ticket: Ticket
}>()

const queryClient = useQueryClient()

const message = ref('')
const time = ref(new Date())
const newTimelineItem = ref(false)

const addCommentMutation = useMutation({
  mutationFn: (): Promise<TimelineItem> =>
    pb.collection('timeline').create({
      ticket: props.ticket.id,
      message: message.value,
      time: time.value
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    message.value = ''
    newTimelineItem.value = false
  },
  onError: handleError
})

const addComment = () => addCommentMutation.mutate()
</script>

<template>
  <Button v-if="!newTimelineItem" variant="outline" @click="newTimelineItem = true">
    <Plus class="mr-2 size-4" />
    Add Timeline Item
  </Button>
  <Card class="p-4" v-else>
    <Popover>
      <PopoverTrigger as-child>
        <Button
          variant="outline"
          :class="cn('justify-start text-left font-normal', !time && 'text-muted-foreground')"
        >
          <CalendarIcon class="mr-2 h-4 w-4" />
          {{ time ? format(time, 'yyyy-MM-dd HH:mm:ss') : 'Pick a date' }}
        </Button>
      </PopoverTrigger>
      <PopoverContent class="w-auto p-0">
        <Calendar v-model="time" mode="datetime" />
      </PopoverContent>
    </Popover>
    <MDEditor
      v-model="message"
      @save="addComment"
      @cancel="newTimelineItem = false"
      placeholder="Add a timeline entry..."
    />
  </Card>
</template>
