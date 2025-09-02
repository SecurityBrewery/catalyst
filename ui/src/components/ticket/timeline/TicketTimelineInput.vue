<script setup lang="ts">
import MDEditor from '@/components/input/MDEditor.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useToast } from '@/components/ui/toast/use-toast'
import { Calendar } from '@/components/ui/v-calendar'

import { Calendar as CalendarIcon, Plus } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { Ticket, TimelineEntry } from '@/client/models'
import { cn, handleError } from '@/lib/utils'

const api = useAPI()

const props = defineProps<{
  ticket: Ticket
}>()

const queryClient = useQueryClient()
const { toast } = useToast()

const message = ref('')
const time = ref(new Date())
const newTimelineItem = ref(false)

const addCommentMutation = useMutation({
  mutationFn: (): Promise<TimelineEntry> =>
    api.createTimeline({
      newTimelineEntry: {
        ticket: props.ticket.id,
        message: message.value,
        time: time.value
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['timeline', props.ticket.id] })
    toast({
      title: 'Timeline item added',
      description: 'The item has been added successfully'
    })
    message.value = ''
    newTimelineItem.value = false
  },
  onError: handleError('Failed to add timeline item')
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
