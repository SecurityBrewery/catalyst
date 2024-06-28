<script lang="ts" setup>
import * as icons from 'lucide-vue-next'

import { computed } from 'vue'

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: Number,
  color: String,
  strokeWidth: Number,
  defaultClass: String
})

const icon = computed((): string => {
  if (!icons[props.name]) {
    return 'Flame'
  }

  return icons[props.name]
})
</script>

<template>
  <component
    :is="icon"
    :size="size"
    :color="color"
    :stroke-width="strokeWidth"
    :default-class="defaultClass"
  />
</template>
