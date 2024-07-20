<script setup lang="ts">
import DynamicMDEditor from '@/components/input/DynamicMDEditor.vue'
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

import { Edit, MoreVertical, Trash } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Comment } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  comment: Comment
}>()

const isOpen = ref(false)
const editMode = ref(false)
const message = ref(props.comment.message)

const deleteCommentMutation = useMutation({
  mutationFn: () => pb.collection('comments').delete(props.comment.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.comment.ticket] })
    isOpen.value = false
  },
  onError: handleError
})

const editCommentMutation = useMutation({
  mutationFn: () =>
    pb.collection('comments').update(props.comment.id, {
      message: message.value
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.comment.ticket] })
    editMode.value = false
  },
  onError: handleError
})

const edit = () => (editMode.value = true)
const save = () => editCommentMutation.mutate()
</script>

<template>
  <div class="rounded-lg border bg-card p-4 text-card-foreground">
    <div class="flex items-start justify-between">
      <div class="flex flex-col gap-1 text-sm">
        <div class="font-semibold">
          {{ comment.expand.author.name }}
        </div>
        <div class="text-xs text-muted-foreground">
          {{ format(new Date(comment.created), 'PPpp') }}
        </div>
      </div>
      <Dialog v-model:open="isOpen">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button size="icon" variant="outline" class="h-8 w-8">
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
              <DialogTitle>Delete comment</DialogTitle>
              <DialogDescription> Are you sure you want to delete this comment?</DialogDescription>
            </DialogHeader>
            <DialogFooter class="sm:justify-start">
              <Button @click="deleteCommentMutation.mutate" variant="destructive">Delete</Button>
              <DialogClose as-child>
                <Button type="button" variant="secondary">Cancel</Button>
              </DialogClose>
            </DialogFooter>
          </DialogContent>
        </DropdownMenu>
      </Dialog>
    </div>
    <DynamicMDEditor v-model="message" v-model:edit="editMode" @save="save" autofocus />
  </div>
</template>
