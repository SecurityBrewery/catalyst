<script setup lang="ts">
import { type HTMLAttributes } from 'vue'

import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div
    :class="
      cn('flex flex-1 items-start justify-start overflow-y-auto overflow-x-hidden', props.class)
    "
  >
    <slot />
  </div>
</template>
