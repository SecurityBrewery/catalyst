<script setup lang="ts">
import DynamicMDEditor from '@/components/input/DynamicMDEditor.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useToast } from '@/components/ui/toast/use-toast'
import { Calendar } from '@/components/ui/v-calendar'

import { Edit, MoreVertical, Trash } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { TimelineEntry } from '@/client/models'
import { cn, handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  timelineItem: TimelineEntry
}>()

const isOpen = ref(false)
const time = ref(new Date(props.timelineItem.time))
const editMode = ref(false)
const message = ref(props.timelineItem.message)

const updateTimelineMutation = useMutation({
  mutationFn: (update: { time?: Date; message?: string }): Promise<TimelineEntry> =>
    api.updateTimeline({
      id: props.timelineItem.id,
      timelineEntryUpdate: {
        time: update.time?.toISOString(),
        message: update.message
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['timeline', props.timelineItem.ticket] })
    toast({
      title: 'Timeline updated',
      description: 'The item has been updated successfully'
    })
    editMode.value = false
  },
  onError: handleError('Failed to update timeline item')
})

watch(
  () => time.value,
  () => {
    if (time.value) {
      updateTimelineMutation.mutate({ time: time.value })
    }
  }
)

const deleteTimelineItemMutation = useMutation({
  mutationFn: () => api.deleteTimeline({ id: props.timelineItem.id }),
  onSuccess: () => {
    queryClient.removeQueries({ queryKey: ['timeline', props.timelineItem.id] })
    queryClient.invalidateQueries({ queryKey: ['timeline', props.timelineItem.ticket] })
    toast({
      title: 'Timeline item deleted',
      description: 'The item has been deleted successfully'
    })
    isOpen.value = false
  },
  onError: handleError('Failed to delete timeline item')
})

const edit = () => (editMode.value = true)
const save = () =>
  updateTimelineMutation.mutate({
    time: time.value,
    message: message.value
  })
</script>

<template>
  <PanelListElement class="flex h-auto flex-row items-start gap-2 p-1 align-top">
    <Popover>
      <PopoverTrigger as-child>
        <Button
          variant="ghost"
          :class="cn('w-20 justify-center text-left font-normal', !time && 'text-muted-foreground')"
        >
          {{ time ? format(time, 'HH:mm:ss') : 'Pick a date' }}
        </Button>
      </PopoverTrigger>
      <PopoverContent class="w-auto p-0">
        <Calendar v-model="time" mode="datetime" />
      </PopoverContent>
    </Popover>
    <div class="my-0.5 flex-1">
      <DynamicMDEditor v-model="message" v-model:edit="editMode" @save="save" autofocus />
    </div>
    <Dialog v-model:open="isOpen">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button size="icon" variant="outline" class="mt-1 mr-1 ml-auto h-8 w-8">
            <MoreVertical class="h-3.5 w-3.5" />
            <span class="sr-only">More</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem class="cursor-pointer" @click="edit">
            <Edit class="mr-2 h-4 w-4" />
            Edit
          </DropdownMenuItem>
          <DropdownMenuItem class="cursor-pointer" as-child>
            <DialogTrigger class="w-full">
              <Trash class="mr-2 h-4 w-4" />
              Delete
            </DialogTrigger>
          </DropdownMenuItem>
        </DropdownMenuContent>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete timeline item</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this timeline item?</DialogDescription
            >
          </DialogHeader>
          <DialogFooter class="sm:justify-start">
            <Button @click="deleteTimelineItemMutation.mutate" variant="destructive">
              Delete
            </Button>
            <DialogClose as-child>
              <Button type="button" variant="secondary"> Cancel</Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </DropdownMenu>
    </Dialog>
  </PanelListElement>
</template>
