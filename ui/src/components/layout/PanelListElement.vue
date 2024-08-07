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
      cn(
        'flex w-full flex-col items-start border-t px-2 py-1 first:rounded-t first:border-none last:rounded-b sm:flex-row sm:items-center',
        props.class
      )
    "
  >
    <slot />
  </div>
</template>
