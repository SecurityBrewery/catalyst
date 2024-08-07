<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { Button } from '@/components/ui/button'

import { ChevronLeft } from 'lucide-vue-next'

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
  <ColumnHeader>
    <Button @click="router.push({ name: 'reactions' })" variant="outline" class="sm:hidden">
      <ChevronLeft class="mr-2 size-4" />
      Back
    </Button>
  </ColumnHeader>

  <ColumnBody>
    <ColumnBodyContainer small>
      <ReactionForm @submit="addReactionMutation.mutate" />
    </ColumnBodyContainer>
  </ColumnBody>
</template>
