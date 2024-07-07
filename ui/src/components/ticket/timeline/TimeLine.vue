<script setup lang="ts">
import { Circle } from 'lucide-vue-next'

import format from 'date-fns/format'

import type { Comment } from '@/lib/types'

defineProps<{
  comments: Array<Comment>
}>()
</script>

<template>
  <ul class="mb-4 [&>li]:grid-cols-[0_min-content_1fr]">
    <li v-for="comment in comments" :key="comment.id" class="grid items-center gap-x-2">
      <!-- item -->
      <p class="col-start-3 col-end-4 mr-auto text-left text-sm text-muted-foreground">
        <!-- heading -->
        {{ format(new Date(comment.created), 'PPpp') }}
      </p>
      <div
        class="col-start-2 col-end-3 row-start-1 row-end-1 flex size-4 items-center justify-center rounded-full border border-current [&>*]:hidden"
      >
        <!-- dot -->
        <Circle className="size-2.5" />
      </div>
      <Separator
        orientation="vertical"
        class="col-start-2 col-end-3 row-start-2 row-end-2 mx-auto flex h-full min-h-16 w-0.5 justify-center rounded-full bg-muted"
      />
      <div class="col-start-3 col-end-4 py-4 text-left">
        <!-- content -->
        {{ comment.message }}
      </div>
    </li>
    <li class="grid items-center gap-x-2">
      <!-- item -->
      <p class="col-start-3 col-end-4 mr-auto text-left text-sm text-muted-foreground">
        <!-- heading -->
        End
      </p>
      <div
        class="col-start-2 col-end-3 row-start-1 row-end-1 flex size-4 items-center justify-center rounded-full border border-current [&>*]:hidden"
      >
        <!-- dot -->
        <Circle className="size-2.5" />
      </div>
    </li>
  </ul>
</template>

<style scoped></style>
