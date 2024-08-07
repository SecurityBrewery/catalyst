<script setup lang="ts">
import { Separator } from '@/components/ui/separator'

import { cn } from '@/lib/utils'

defineProps<{
  title?: string
  nowrap?: boolean
  hideSeparator?: boolean
}>()
</script>

<template>
  <div
    :class="
      cn('flex min-h-14 flex-wrap items-center gap-2 bg-background p-2', nowrap && 'flex-nowrap')
    "
  >
    <h1 v-if="title" class="text-xl font-bold">
      {{ title }}
    </h1>
    <slot />
  </div>
  <Separator v-if="!hideSeparator" />
</template>
