<script setup lang="ts">
import { Circle, CircleCheck, CircleMinus, CircleX } from 'lucide-vue-next'

defineProps<{
  status: 'open' | 'completed' | 'pending' | 'failed'
}>()

const strokeWidth = 2
</script>

<template>
  <Circle v-if="status === 'open'" :stroke-width="strokeWidth" class="ml-2 h-4 w-4 text-gray-500" />
  <CircleCheck
    v-else-if="status === 'completed'"
    :stroke-width="strokeWidth"
    class="ml-2 h-4 w-4 text-emerald-500"
  />
  <CircleMinus
    v-else-if="status === 'pending'"
    :stroke-width="strokeWidth"
    class="ml-2 h-4 w-4 text-amber-500"
  />
  <CircleX
    v-else-if="status === 'failed'"
    :stroke-width="strokeWidth"
    class="ml-2 h-4 w-4 text-rose-500"
  />
</template>
