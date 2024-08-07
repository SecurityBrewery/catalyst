<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { Button } from '@/components/ui/button'
import { toast } from '@/components/ui/toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Reaction } from '@/lib/types'
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
  queryFn: (): Promise<Reaction> => pb.collection('reactions').getOne(props.id)
})

const updateReactionMutation = useMutation({
  mutationFn: (update: any) => pb.collection('reactions').update(props.id, update),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['reactions'] }),
  onError: handleError
})

onMounted(() => {
  if (props.id) {
    pb.collection('reactions').subscribe(props.id, (data) => {
      if (data.action === 'delete') {
        toast({
          title: 'Reaction deleted',
          description: 'The reaction has been deleted.',
          variant: 'destructive'
        })

        router.push({ name: 'reactions' })

        return
      }

      if (data.action === 'update') {
        toast({
          title: 'Reaction updated',
          description: 'The reaction has been updated.'
        })

        queryClient.invalidateQueries({ queryKey: ['reactions', props.id] })
      }
    })
  }
})

onUnmounted(() => {
  if (props.id) {
    pb.collection('reactions').unsubscribe(props.id)
  }
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
          collection="reactions"
          :id="reaction.id"
          :name="reaction.name"
          :singular="'Reaction'"
          :to="{ name: 'reactions' }"
          :queryKey="['reactions']"
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
