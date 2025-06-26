<script setup lang="ts">
import Icon from '@/components/Icon.vue'

import { formatDistanceToNow } from 'date-fns'

import { cn } from '@/lib/utils'

defineProps<{
  icon?: string
  title: string
  subtitle: string
  description: string
  created: Date

  open: boolean
  active: boolean
  to: string | { name: string; params: Record<string, string | number> }
}>()
</script>

<template>
  <RouterLink
    :class="
      cn(
        'bg-card hover:bg-accent flex flex-col items-start gap-2 rounded-lg border p-3 text-left text-sm transition-all',
        active && 'bg-accent'
      )
    "
    :to="to"
  >
    <div class="flex w-full flex-col gap-1">
      <div class="flex items-center">
        <div class="flex items-center gap-1">
          <Icon v-if="icon" :name="icon" class="mr-2 size-4" />
          <div class="font-semibold">
            {{ title }}
          </div>
          <span v-if="open" class="ml-1 flex h-2 w-2 rounded-full bg-blue-600" />
        </div>
        <div :class="cn('ml-auto text-xs', active ? 'text-foreground' : 'text-muted-foreground')">
          {{ formatDistanceToNow(created, { addSuffix: true }) }}
        </div>
      </div>

      <div v-if="subtitle" class="text-xs font-medium">
        {{ subtitle }}
      </div>
    </div>
    <div v-if="description" class="text-muted-foreground line-clamp-2 text-xs">
      {{ description }}
    </div>
  </RouterLink>
</template>
