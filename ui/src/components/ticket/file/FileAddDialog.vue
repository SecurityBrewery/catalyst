<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { File as CFile, Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const isOpen = defineModel<boolean>()
const file = ref<File | null>(null)

const addFileMutation = useMutation({
  mutationFn: (): Promise<CFile> => {
    if (!file.value) return Promise.reject('No file selected')

    return pb.collection('files').create({
      ticket: props.ticket.id,
      name: file.value.name,
      blob: file.value,
      size: file.value.size
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    isOpen.value = false
  },
  onError: handleError
})

const save = () => addFileMutation.mutate()

function handleFileUpload($event: Event) {
  const target = $event.target as HTMLInputElement
  if (target && target.files) {
    file.value = target.files[0]
  }
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New File</DialogTitle>
        <DialogDescription> Upload a new file to this ticket.</DialogDescription>
      </DialogHeader>

      <Input id="file" type="file" class="mt-2" @change="handleFileUpload($event)" />

      <DialogFooter class="mt-2 sm:justify-start">
        <Button @click="save">Upload</Button>
        <DialogClose as-child>
          <Button type="button" variant="secondary">Cancel</Button>
        </DialogClose>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
