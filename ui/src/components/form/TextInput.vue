<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import { Save, X } from 'lucide-vue-next'

import { type HTMLAttributes, ref, watchEffect } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    placeholder?: string
    class?: HTMLAttributes['class']
  }>(),
  {
    modelValue: '',
    placeholder: ''
  }
)

const emit = defineEmits(['update:modelValue'])

const text = ref(props.modelValue)
const editMode = ref(false)
const input = ref<HTMLInputElement | null>(null)

const setValue = () => emit('update:modelValue', text.value)

watchEffect(() => {
  if (editMode.value && input.value) {
    input.value.$el.focus()
  }
})

const edit = () => {
  text.value = props.modelValue
  editMode.value = true
}

const cancel = () => {
  text.value = props.modelValue
  editMode.value = false
}
</script>

<template>
  <Button v-if="!editMode" variant="outline" size="icon" @click="edit" class="flex-1">
    <div class="ml-3 w-full text-start font-normal">
      {{ text }}
    </div>
  </Button>
  <div v-else class="flex w-full flex-row gap-2">
    <Input
      ref="input"
      v-model="text"
      :placeholder="placeholder"
      @keydown.enter="setValue"
      class="flex-1"
    />
    <Button variant="outline" size="icon" @click="cancel" class="shrink-0">
      <X class="size-4" />
      <span class="sr-only">Cancel</span>
    </Button>
    <Button
      variant="outline"
      size="icon"
      @click="setValue"
      :disabled="text === modelValue"
      class="shrink-0"
    >
      <Save class="size-4" />
      <span class="sr-only">Save</span>
    </Button>
  </div>
</template>
