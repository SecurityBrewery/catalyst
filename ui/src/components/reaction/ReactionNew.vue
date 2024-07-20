<script setup lang="ts">
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Reaction, Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()
const router = useRouter()

const addReactionMutation = useMutation({
  mutationFn: (values: Reaction): Promise<Reaction> => pb.collection('reactions').create(values),
  onSuccess: (data: Ticket) => {
    router.push({ name: 'reactions', params: { id: data.id } })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
  },
  onError: handleError
})
</script>

<template>
  <div class="flex h-full flex-1 flex-col overflow-hidden">
    <div class="flex min-h-14 items-center bg-background px-4 py-2"></div>
    <Separator />

    <ScrollArea class="flex-1">
      <div class="flex max-w-[640px] flex-col gap-4 p-4">
        <ReactionForm @submit="addReactionMutation.mutate" />
      </div>
    </ScrollArea>
  </div>
</template>
