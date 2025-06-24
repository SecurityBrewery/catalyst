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
import { useToast } from '@/components/ui/toast/use-toast'

import { useQueryClient } from '@tanstack/vue-query'
import Uppy from '@uppy/core'
import Tus from '@uppy/tus'
import { Dashboard } from '@uppy/vue'

import type { Ticket } from '@/client/models'
import { useAuthStore } from '@/store/auth'

const queryClient = useQueryClient()
const { toast } = useToast()

const authStore = useAuthStore()

const props = defineProps<{
  ticket: Ticket
}>()

const isOpen = defineModel<boolean>()

const uppy = new Uppy().use(Tus, {
  endpoint: '/files/',
  headers: {
    'X-Ticket-ID': props.ticket.id,
    authorization: `Bearer ${authStore.token}`
  },
  onSuccess() {
    queryClient.invalidateQueries({ queryKey: ['files', props.ticket.id] })
    toast({
      title: 'File uploaded',
      description: 'The file has been uploaded successfully'
    })
  }
})
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New File</DialogTitle>
        <DialogDescription> Upload a new file to this ticket.</DialogDescription>
      </DialogHeader>

      <Dashboard :uppy="uppy" :props="{ width: '100%', height: '350px' }" />

      <DialogFooter class="mt-2">
        <DialogClose as-child>
          <Button type="button" variant="secondary">Close</Button>
        </DialogClose>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
