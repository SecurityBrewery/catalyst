<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  size: number
}>()
</script>

<template>
  <img src="@/assets/flask.svg" alt="Catalyst Logo" :class="cn('dark:hidden', 'size-' + size)" />
  <img
    src="@/assets/flask_white.svg"
    alt="Catalyst Logo"
    :class="cn('hidden dark:flex', 'size-' + size)"
  />
</template>
