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
import type { File, Ticket } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  file: File
}>()

const isOpen = ref(false)

const removeFileMutation = useMutation({
  mutationFn: (): Promise<boolean> => pb.collection('files').delete(props.file.id),
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
      <Button variant="ghost" size="icon">
        <Trash2 class="h-4 w-4" />
      </Button>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle> Delete File "{{ props.file.name }}" </DialogTitle>
        <DialogDescription> Are you sure you want to delete this file? </DialogDescription>
      </DialogHeader>

      <DialogFooter class="mt-2">
        <DialogClose as-child>
          <Button type="button" variant="secondary"> Cancel</Button>
        </DialogClose>
        <Button type="button" variant="destructive" @click="removeFileMutation.mutate()">
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
