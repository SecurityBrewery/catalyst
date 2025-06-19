<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <tfoot :class="cn('bg-muted/50 border-t font-medium last:[&>tr]:border-b-0', props.class)">
    <slot />
  </tfoot>
</template>
