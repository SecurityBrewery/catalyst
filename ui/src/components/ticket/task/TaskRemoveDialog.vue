<script setup lang="ts">
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
import { toast } from '@/components/ui/toast'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Task, Ticket } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  task: Task
}>()

const isOpen = ref(false)

const removeTaskMutation = useMutation({
  mutationFn: (): Promise<boolean> => pb.collection('tasks').delete(props.task.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    isOpen.value = false
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="ghost" size="icon" class="h-8 w-8">
        <Trash2 class="size-4" />
      </Button>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle> Delete Task "{{ props.task.name }}" </DialogTitle>
        <DialogDescription> Are you sure you want to delete this task? </DialogDescription>
      </DialogHeader>

      <DialogFooter class="mt-2">
        <DialogClose as-child>
          <Button type="button" variant="secondary"> Cancel</Button>
        </DialogClose>
        <Button type="button" variant="destructive" @click="removeTaskMutation.mutate()">
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
