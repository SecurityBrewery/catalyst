<script setup lang="ts">
import { formatDistanceToNow } from 'date-fns'

import { cn } from '@/lib/utils'

defineProps<{
  title: string
  subtitle: string
  description: string
  created: string

  open: boolean
  active: boolean
  to: string | { name: string; params: Record<string, string | number> }
}>()
</script>

<template>
  <RouterLink
    :class="
      cn(
        'flex flex-col items-start gap-2 rounded-lg border bg-card p-3 text-left text-sm transition-all hover:bg-accent',
        active && 'bg-accent'
      )
    "
    :to="to"
  >
    <div class="flex w-full flex-col gap-1">
      <div class="flex items-center">
        <div class="flex items-center gap-2">
          <div class="font-semibold">
            {{ title }}
          </div>
          <span v-if="open" class="flex h-2 w-2 rounded-full bg-blue-600" />
        </div>
        <div :class="cn('ml-auto text-xs', active ? 'text-foreground' : 'text-muted-foreground')">
          {{ formatDistanceToNow(new Date(created), { addSuffix: true }) }}
        </div>
      </div>

      <div v-if="subtitle" class="text-xs font-medium">
        {{ subtitle }}
      </div>
    </div>
    <div v-if="description" class="line-clamp-2 text-xs text-muted-foreground">
      {{ description }}
    </div>
  </RouterLink>
</template>
