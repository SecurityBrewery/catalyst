<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ReactionWebhookForm from '@/components/reaction/ReactionWebhookForm.vue'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { ReactionWebhook } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const {
  isPending,
  isError,
  data: reaction,
  error
} = useQuery({
  queryKey: ['reactions_webhooks', props.id],
  queryFn: (): Promise<ReactionWebhook> => pb.collection('reactions_webhooks').getOne(props.id)
})

const updateReactionMutation = useMutation({
  mutationFn: (update: any) => pb.collection('reactions_webhooks').update(props.id, update),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    queryClient.invalidateQueries({ queryKey: ['reactions_webhooks'] })
  },
  onError: handleError
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error" :value="reaction">
    <div class="flex h-full flex-1 flex-col overflow-hidden">
      <div class="flex items-center bg-background px-4 py-2">
        <div class="ml-auto">
          <DeleteDialog
            v-if="reaction"
            collection="reactions_webhooks"
            :id="reaction.id"
            :name="reaction.name"
            :singular="'Reaction'"
            :to="{ name: 'reactions' }"
            :queryKey="['reactions']"
          />
        </div>
      </div>
      <Separator />

      <ScrollArea v-if="reaction" class="flex-1">
        <div class="flex max-w-[640px] flex-col gap-4 p-4">
          <h1 class="text-3xl font-bold">Webhook Reaction Editor</h1>

          <p class="py-4 text-sm text-muted-foreground">
            Sent a POST request to the specified destination URL with the specified headers.
          </p>

          <ReactionWebhookForm :reaction="reaction" @submit="updateReactionMutation.mutate" />
        </div>
      </ScrollArea>
    </div>
  </TanView>
</template>
