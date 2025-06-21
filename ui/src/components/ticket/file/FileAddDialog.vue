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
import { useToast } from '@/components/ui/toast/use-toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { ModelFile, Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  ticket: Ticket
}>()

const isOpen = defineModel<boolean>()
const file = ref<File | null>(null)

const addFileMutation = useMutation({
  mutationFn: async (): Promise<ModelFile> => {
    if (!file.value) return Promise.reject('No file selected')

    return api.createFile({
      newFile: {
        ticket: props.ticket.id,
        name: file.value.name,
        blob: (await toBase64(file.value)) as string,
        size: file.value.size
      }
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['files', props.ticket.id] })
    toast({
      title: 'File uploaded',
      description: 'The file has been uploaded successfully'
    })
    isOpen.value = false
  },
  onError: handleError('Failed to upload file')
})

const toBase64 = (file: File) =>
  new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
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
