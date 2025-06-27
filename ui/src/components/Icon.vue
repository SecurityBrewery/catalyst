<script lang="ts" setup>
import * as icons from 'lucide-vue-next'

import { computed } from 'vue'

const props = defineProps({
  name: String,
  size: Number,
  color: String,
  strokeWidth: Number,
  defaultClass: String
})

const IconComponent = computed(() => {
  if (!props.name) {
    return icons.HelpCircle
  }

  const component = icons[props.name]
  if (!component) {
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
