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
        'flex w-full items-center border-t px-2 py-1 first:rounded-t first:border-none last:rounded-b',
        props.class
      )
    "
  >
    <slot />
  </div>
</template>
