<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  small?: boolean
}>()
</script>

<template>
  <div :class="cn('mx-auto flex w-full max-w-6xl flex-col gap-4 p-4', small && 'max-w-188')">
    <slot />
  </div>
</template>
