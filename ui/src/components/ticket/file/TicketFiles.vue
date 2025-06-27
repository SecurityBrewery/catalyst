<script setup lang="ts">
import '@uppy/core/dist/style.min.css'
import '@uppy/dashboard/dist/style.min.css'

import DeleteDialog from '@/components/common/DeleteDialog.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import FileAddDialog from '@/components/ticket/file/FileAddDialog.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { Download, Trash2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { ModelFile, Ticket } from '@/client/models'
import { handleError, human } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const authStore = useAuthStore()

const props = defineProps<{
  ticket: Ticket
  files: Array<ModelFile> | undefined
}>()

const downloadFile = (file: any) => {
  fetch(`/api/files/${file.id}/download`, {
    headers: { Authorization: `Bearer ${authStore.token}` }
  })
    .then((response) => response.blob())
    .then((blob) => {
      var _url = window.URL.createObjectURL(blob)

      // Create a link element, set its href to the blob URL, and trigger a download
      const link = document.createElement('a')
      link.href = _url
      link.download = file.name
      document.body.appendChild(link)
      link.click()
    })
    .catch((err) => {
      console.log(err)
    })
}

const dialogOpen = ref(false)

const isDemo = ref(false)

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: () => api.getConfig()
})

const deleteMutation = useMutation({
  mutationFn: (id: string) => api.deleteFile({ id }),
  onSuccess: (data, id) => {
    queryClient.removeQueries({ queryKey: ['files', id] })
    queryClient.invalidateQueries({ queryKey: ['files', props.ticket.id] })
    toast({
      title: 'File deleted',
      description: 'The file has been deleted successfully'
    })
  },
  onError: handleError('Failed to delete file')
})

watch(
  () => config.value,
  (newConfig) => {
    if (!newConfig) return
    if (newConfig.flags.includes('demo')) {
      isDemo.value = true
    }
  },
  { immediate: true }
)
</script>

<template>
  <TicketPanel title="Files" @add="dialogOpen = true" :hideAdd="isDemo">
    <FileAddDialog v-if="!isDemo" v-model="dialogOpen" :ticket="ticket" />
    <div
      v-if="!files || files.length === 0"
      class="flex h-10 items-center p-4 text-muted-foreground"
    >
      {{ isDemo ? 'Cannot upload files in demo mode' : 'No files added yet.' }}
    </div>
    <div
      v-for="file in files"
      :key="file.id"
      :title="file.name"
      class="flex w-full items-center border-t py-1 pl-2 pr-1 first:rounded-t first:border-none last:rounded-b"
    >
      <div class="flex flex-1 items-center overflow-hidden pr-2">
        {{ file.name }}

        <div class="ml-1 flex-1 text-nowrap text-sm text-muted-foreground">
          {{ human(file.size) }}
        </div>
      </div>

      <Button
        variant="ghost"
        size="icon"
        class="mr-1 size-8 text-muted-foreground"
        @click="downloadFile(file)"
      >
        <Download class="size-4" />
      </Button>
      <DeleteDialog
        v-if="file"
        :name="file.name"
        singular="File"
        @delete="deleteMutation.mutate(file.id)"
      >
        <Button variant="ghost" size="icon" class="size-8">
          <Trash2 class="size-4" />
          <span class="sr-only">Delete File</span>
        </Button>
      </DeleteDialog>
    </div>
  </TicketPanel>
</template>
