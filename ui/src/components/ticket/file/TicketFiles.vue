<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import FileAddDialog from '@/components/ticket/file/FileAddDialog.vue'
import { Button } from '@/components/ui/button'

import { Download, Trash2 } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { ModelFile, Ticket } from '@/client/models'
import { human } from '@/lib/utils'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  files: Array<ModelFile> | undefined
}>()

const downloadFile = (file: any) => {
  window.open(`/api/files/files/${file.id}/${file.blob}?download=1`, '_blank')
}

const dialogOpen = ref(false)

const isDemo = ref(false)

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: (): Promise<Record<string, Array<String>>> => {
    return fetch('/config').then((response) => {
      if (response.ok) {
        return response.json()
      }

      throw new Error('Failed to fetch config')
    })
  }
})

const deleteMutation = useMutation({
  mutationFn: (id: string) => {
    return api.deleteFile({ id })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['files', props.ticket.id] })
  },
  onError: handleError
})

watch(
  () => config.value,
  (newConfig) => {
    if (!newConfig) return
    if (newConfig['flags'].includes('demo')) {
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
      class="flex w-full items-center border-t first:rounded-t first:border-none last:rounded-b"
    >
      <div class="flex flex-1 items-center overflow-hidden py-2 pl-4 pr-2">
        {{ file.name }}

        <div class="ml-1 flex-1 text-nowrap text-sm text-muted-foreground">
          {{ human(file.size) }}
        </div>
      </div>

      <Button
        variant="ghost"
        size="icon"
        class="mr-1 text-muted-foreground"
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
        <Button variant="ghost" size="icon" class="h-8 w-8">
          <Trash2 class="size-4" />
        </Button>
      </DeleteDialog>
    </div>
  </TicketPanel>
</template>
