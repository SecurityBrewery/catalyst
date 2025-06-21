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
import { useToast } from '@/components/ui/toast/use-toast'

import { Edit, MoreVertical, Trash } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import format from 'date-fns/format'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { ExtendedComment } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  comment: ExtendedComment
}>()

const isOpen = ref(false)
const editMode = ref(false)
const message = ref(props.comment.message)

const deleteCommentMutation = useMutation({
  mutationFn: () => api.deleteComment({ id: props.comment.id }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['comments', props.comment.ticket] })
    toast({
      title: 'Comment deleted',
      description: 'The comment has been deleted successfully'
    })
    isOpen.value = false
  },
  onError: handleError('Failed to delete comment')
})

const editCommentMutation = useMutation({
  mutationFn: () =>
    api.updateComment({
      id: props.comment.id,
      commentUpdate: {
        message: message.value
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['comments', props.comment.ticket] })
    toast({
      title: 'Comment updated',
      description: 'The comment has been updated successfully'
    })
    editMode.value = false
  },
  onError: handleError('Failed to update comment')
})

const edit = () => (editMode.value = true)
const save = () => editCommentMutation.mutate()
</script>

<template>
  <div class="bg-card text-card-foreground rounded-lg border p-4">
    <div class="flex items-start justify-between">
      <div class="flex flex-col gap-1 text-sm">
        <div class="font-semibold">
          {{ comment.authorName }}
        </div>
        <div class="text-muted-foreground text-xs">
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
