<script setup lang="ts">
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import FileAddDialog from '@/components/ticket/file/FileAddDialog.vue'
import FileRemoveDialog from '@/components/ticket/file/FileRemoveDialog.vue'
import { Button } from '@/components/ui/button'

import { Download } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { ref, watch } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { File, Ticket } from '@/lib/types'
import { human } from '@/lib/utils'

defineProps<{
  ticket: Ticket
  files: Array<File> | undefined
}>()

const downloadFile = (file: any) => {
  window.open(`/api/files/files/${file.id}/${file.blob}?download=1`, '_blank')
}

const dialogOpen = ref(false)

const isDemo = ref(false)

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: (): Promise<Record<string, Array<String>>> => pb.send('/api/config', {})
})

watch(config, (newConfig) => {
  if (!newConfig) return
  if (newConfig['flags'].includes('demo')) {
    isDemo.value = true
  }
})
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
      <FileRemoveDialog :ticket="ticket" :file="file" />
    </div>
  </TicketPanel>
</template>