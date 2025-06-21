<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Reaction } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

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
  onSuccess: () => {
    toast({
      title: 'Reaction updated',
      description: 'The reaction has been updated successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
  },
  onError: handleError('Failed to update reaction')
})

const deleteMutation = useMutation({
  mutationFn: () => api.deleteReaction({ id: props.id }),
  onSuccess: () => {
    queryClient.removeQueries({ queryKey: ['reactions', props.id] })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    toast({
      title: 'Reaction deleted',
      description: 'The reaction has been deleted successfully'
    })
    router.push({ name: 'reactions' })
  },
  onError: handleError('Failed to delete reaction')
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
