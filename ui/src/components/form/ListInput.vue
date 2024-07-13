<script setup lang="ts">
import TextInput from '@/components/form/TextInput.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import { Plus, Trash2 } from 'lucide-vue-next'

import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string[]
    placeholder?: string
  }>(),
  {
    modelValue: () => [],
    placeholder: ''
  }
)

const emit = defineEmits(['update:modelValue'])

const newItem = ref('')

const updateModelValue = (value: string, index: number) => {
  const newValue = props.modelValue
  newValue[index] = value
  emit('update:modelValue', newValue)
}

const addModelValue = () => {
  emit('update:modelValue', [...props.modelValue, newItem.value])
  newItem.value = ''
}

const removeModelValue = (index: number) =>
  emit(
    'update:modelValue',
    props.modelValue.filter((_, i) => i !== index)
  )
</script>

<template>
  <div class="flex flex-col gap-2">
    <div v-for="(item, index) in modelValue" :key="item" class="flex flex-row items-center gap-2">
      <TextInput
        :modelValue="item"
        @update:modelValue="updateModelValue($event, index)"
        :placeholder="placeholder"
      />
      <Button variant="outline" size="icon" @click="removeModelValue(index)" class="shrink-0">
        <Trash2 class="size-4" />
        <span class="sr-only">Remove item</span>
      </Button>
    </div>
    <div class="flex flex-row items-center gap-2">
      <Input v-model="newItem" :placeholder="placeholder" @keydown.enter.prevent="addModelValue" />
      <Button
        variant="outline"
        size="icon"
        @click="addModelValue"
        :disabled="newItem === ''"
        class="shrink-0"
      >
        <Plus class="size-4" />
        <span class="sr-only">Add item</span>
      </Button>
    </div>
  </div>
</template>
