<script setup lang="ts">
import TicketComment from '@/components/ticket/comment/TicketComment.vue'
import TicketCommentInput from '@/components/ticket/comment/TicketCommentInput.vue'
import { Card } from '@/components/ui/card'

import type { Comment, Ticket } from '@/lib/types'

defineProps<{
  ticket: Ticket
  comments: Array<Comment> | undefined
}>()
</script>

<template>
  <div class="mt-2 flex flex-col gap-2">
    <Card
      v-if="!comments || comments.length === 0"
      class="flex h-10 items-center p-4 text-muted-foreground"
    >
      No comments added yet.
    </Card>
    <div v-else class="flex flex-col gap-2">
      <TicketComment
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        class="rounded-lg border bg-card p-4 text-card-foreground"
      />
    </div>
    <TicketCommentInput :ticket="ticket" class="w-full" />
  </div>
</template>
