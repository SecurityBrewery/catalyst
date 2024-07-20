<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Reaction } from '@/lib/types'
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
  queryKey: ['reactions', props.id],
  queryFn: (): Promise<Reaction> => pb.collection('reactions').getOne(props.id)
})

const updateReactionMutation = useMutation({
  mutationFn: (update: any) => pb.collection('reactions').update(props.id, update),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['reactions'] }),
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
            collection="reactions"
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
          <ReactionForm :reaction="reaction" @submit="updateReactionMutation.mutate" hide-cancel />
        </div>
      </ScrollArea>
    </div>
  </TanView>
</template>
