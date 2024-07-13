<script setup lang="ts">
import ReactionForm from '@/components/reaction/ReactionForm.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogScrollContent,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Reaction, Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()
const router = useRouter()

const isOpen = ref(false)

const addReactionMutation = useMutation({
  mutationFn: (values: Reaction): Promise<Reaction> => pb.collection('reactions').create(values),
  onSuccess: (data: Ticket) => {
    router.push({ name: 'reactions', params: { id: data.id } })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    isOpen.value = false
  },
  onError: handleError
})
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="ghost">New Reaction</Button>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New Reaction</DialogTitle>
        <DialogDescription>Create a new reaction</DialogDescription>
      </DialogHeader>

      <DialogScrollContent>
        <ReactionForm @submit="addReactionMutation.mutate" />
      </DialogScrollContent>
    </DialogContent>
  </Dialog>
</template>
