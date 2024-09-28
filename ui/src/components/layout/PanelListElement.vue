<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div
    :class="
      cn('flex w-full flex-col items-start px-2 py-1 sm:flex-row sm:items-center', props.class)
    "
  >
    <slot />
  </div>
</template>
