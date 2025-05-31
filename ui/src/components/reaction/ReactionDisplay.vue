<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { Button } from '@/components/ui/button'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { api } from '@/api'
import type { Reaction } from '@/client/models'
import { handleError } from '@/lib/utils'

const router = useRouter()
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
  queryFn: (): Promise<Reaction> => api.getReaction({ id: props.id })
})

const updateReactionMutation = useMutation({
  mutationFn: (update: any) => api.updateReaction({ id: props.id, reactionUpdate: update }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['reactions'] }),
  onError: handleError
})

const deleteMutation = useMutation({
  mutationFn: () => {
    return api.deleteReaction({ id: props.id })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    router.push({ name: 'reactions' })
  },
  onError: handleError
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader>
      <Button @click="router.push({ name: 'reactions' })" variant="outline" class="sm:hidden">
        <ChevronLeft class="mr-2 size-4" />
        Back
      </Button>
      <div class="ml-auto">
        <DeleteDialog
          v-if="reaction"
          :name="reaction.name"
          singular="Reaction"
          @delete="deleteMutation.mutate"
        />
      </div>
    </ColumnHeader>

    <ColumnBody v-if="reaction">
      <ColumnBodyContainer small>
        <ReactionForm :reaction="reaction" @submit="updateReactionMutation.mutate" />
      </ColumnBodyContainer>
    </ColumnBody>
  </TanView>
</template>
