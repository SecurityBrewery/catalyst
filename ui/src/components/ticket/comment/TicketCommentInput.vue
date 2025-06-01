<script setup lang="ts">
import MDEditor from '@/components/input/MDEditor.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

import { Plus } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { Comment, Ticket } from '@/client/models'
import { handleError } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const api = useAPI()

const authStore = useAuthStore()

const props = defineProps<{
  ticket: Ticket
}>()

const queryClient = useQueryClient()

const message = ref('')
const isOpen = ref(false)

const addCommentMutation = useMutation({
  mutationFn: (): Promise<Comment> => {
    if (!authStore.user) return Promise.reject('Not authenticated')
    return api.createComment({
      newComment: {
        ticket: props.ticket.id,
        author: authStore.user.id,
        message: message.value
      }
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
