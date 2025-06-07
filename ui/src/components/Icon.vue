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

const IconComponent = computed(() => {
  const component = icons[props.name]
  if (!component) {
    console.warn(`Icon "${props.name}" not found`)
    return icons.HelpCircle
  }
  return component
})
</script>

<template>
  <component
    :is="IconComponent"
    :size="size"
    :color="color"
    :stroke-width="strokeWidth"
    :class="defaultClass"
  />
</template>
