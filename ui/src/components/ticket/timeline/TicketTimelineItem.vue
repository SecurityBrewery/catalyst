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
import { Calendar } from '@/components/ui/v-calendar'

import { Edit, MoreVertical, Trash } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref, watch } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { TimelineItem } from '@/lib/types'
import { cn, handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  timelineItem: TimelineItem
}>()

const isOpen = ref(false)
const time = ref(props.timelineItem.time)
const editMode = ref(false)
const message = ref(props.timelineItem.message)

const updateTimelineMutation = useMutation({
  mutationFn: (update: { time?: string; message?: string }): Promise<TimelineItem> =>
    pb.collection('timeline').update(props.timelineItem.id, update),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.timelineItem.ticket] })
    editMode.value = false
  },
  onError: handleError
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
  mutationFn: () => pb.collection('timeline').delete(props.timelineItem.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.timelineItem.ticket] })
    isOpen.value = false
  },
  onError: handleError
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
          <Button size="icon" variant="outline" class="ml-auto mr-1 mt-1 h-8 w-8">
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
