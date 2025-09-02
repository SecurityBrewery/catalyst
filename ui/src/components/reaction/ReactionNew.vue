<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Reaction } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()
const { toast } = useToast()

const addReactionMutation = useMutation({
  mutationFn: (values: Reaction): Promise<Reaction> => api.createReaction({ newReaction: values }),
  onSuccess: (data: Reaction) => {
    router.push({ name: 'reactions', params: { id: data.id } })
    toast({
      title: 'Reaction created',
      description: 'The reaction has been created successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
  },
  onError: handleError('Failed to create reaction')
})
</script>

<template>
  <ColumnHeader>
    <Button @click="router.push({ name: 'reactions' })" variant="outline" class="md:hidden">
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
