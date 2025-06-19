<script setup lang="ts">
import { type HTMLAttributes, onMounted, ref } from 'vue'

import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()

const modelValue = defineModel<string>({
  default: ''
})

const textarea = ref<HTMLElement | null>(null)

const resize = () => {
  if (!textarea.value) return
  textarea.value.style.height = 'auto' // Reset to default or minimum height
  textarea.value.style.height = textarea.value.scrollHeight + 2 + 'px'
}

onMounted(() => resize())
</script>

<template>
  <textarea
    ref="textarea"
    rows="1"
    @focus="resize"
    @input="resize"
    v-model="modelValue"
    :class="
      cn(
        'border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-20 w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50',
        props.class
      )
    "
  />
</template>
