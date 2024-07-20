<script setup lang="ts">
import MultiSelect from '@/components/form/MultiSelect.vue'

import { computed } from 'vue'

const modelValue = defineModel<string[]>({
  default: []
})

const items = ['Tickets', 'Tasks', 'Comments', 'Timeline', 'Links', 'Files']

const mapping: Record<string, string> = {
  tickets: 'Tickets',
  tasks: 'Tasks',
  comments: 'Comments',
  timeline: 'Timeline',
  links: 'Links',
  files: 'Files'
}

const niceNames = computed(() => modelValue.value.map((collection) => mapping[collection]))

const updateModelValue = (values: string[]) => {
  modelValue.value = values.map(
    (value) => Object.keys(mapping).find((key) => mapping[key] === value)!
  )
}
</script>

<template>
  <MultiSelect
    :modelValue="niceNames"
    @update:modelValue="updateModelValue"
    :items="items"
    placeholder="Select collections..."
  />
</template>
