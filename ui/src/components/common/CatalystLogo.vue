<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: string
}>()
</script>

<template>
  <img src="@/assets/flask.svg" alt="Catalyst Logo" :class="cn('dark:hidden', props.class)" />
  <img
    src="@/assets/flask_white.svg"
    alt="Catalyst Logo"
    :class="cn('hidden dark:flex', props.class)"
  />
</template>
