<script setup lang="ts">
import MDEditor from '@/components/input/MDEditor.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

import { Plus } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const props = defineProps<{
  ticket: Ticket
}>()

const queryClient = useQueryClient()

const message = ref('')
const isOpen = ref(false)

const addCommentMutation = useMutation({
  mutationFn: (): Promise<Comment> => {
    if (!pb.authStore.model) return Promise.reject('Not authenticated')
    return pb.collection('comments').create({
      ticket: props.ticket.id,
      author: pb.authStore.model.id,
      message: message.value
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    message.value = ''
  },
  onError: handleError
})
</script>

<template>
  <Button v-if="!isOpen" variant="outline" @click="isOpen = true">
    <Plus class="mr-2 size-4" />
    Add Comment
  </Button>
  <Card class="p-4" v-else>
    <MDEditor
      v-model="message"
      @save="addCommentMutation.mutate"
      @cancel="isOpen = false"
      autofocus
      placeholder="Type a comment..."
    />
  </Card>
</template>
