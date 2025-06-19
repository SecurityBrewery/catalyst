<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <h5 :class="cn('mb-1 leading-none font-medium tracking-tight', props.class)">
    <slot />
  </h5>
</template>
