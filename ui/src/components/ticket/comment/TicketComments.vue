<script setup lang="ts">
import TicketComment from '@/components/ticket/comment/TicketComment.vue'
import TicketCommentInput from '@/components/ticket/comment/TicketCommentInput.vue'
import { Card } from '@/components/ui/card'

import type { ExtendedComment, Ticket } from '@/client/models'

defineProps<{
  ticket: Ticket
  comments: Array<ExtendedComment> | undefined
}>()
</script>

<template>
  <div class="mt-2 flex flex-col gap-2">
    <Card
      v-if="!comments || comments.length === 0"
      class="text-muted-foreground flex h-10 items-center p-4"
    >
      No comments added yet.
    </Card>
    <div v-else class="flex flex-col gap-2">
      <TicketComment
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        class="bg-card text-card-foreground rounded-lg border p-4"
      />
    </div>
    <TicketCommentInput :ticket="ticket" class="w-full" />
  </div>
</template>
