<script setup lang="ts">
import MultiSelect from '@/components/form/MultiSelect.vue'

import { computed } from 'vue'

const modelValue = defineModel<string[]>({
  default: []
})

const items = ['Create Events', 'Update Events', 'Delete Events']

const mapping: Record<string, string> = {
  create: 'Create Events',
  update: 'Update Events',
  delete: 'Delete Events'
}

const niceNames = computed(() =>
  modelValue.value.map((collection) => mapping[collection] as string)
)

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
    placeholder="Select events..."
  />
</template>
