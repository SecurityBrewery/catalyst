<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  small?: boolean
}>()
</script>

<template>
  <div :class="cn('mx-auto flex w-full max-w-[72rem] gap-4 p-4', small && 'max-w-[47rem]')">
    <slot />
  </div>
</template>
