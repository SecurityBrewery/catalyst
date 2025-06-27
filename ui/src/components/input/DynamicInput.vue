<script setup lang="ts">
import ShortCut from '@/components/ShortCut.vue'

import { ref } from 'vue'

const model = defineModel({
  type: String
})

const props = defineProps({
  placeholder: {
    type: String,
    required: false
  },
  type: {
    type: String,
    required: false,
    default: 'input',
    validator: (value: string) => ['input', 'textarea'].includes(value)
  }
})

const active = ref(false)
const text = ref(model.value)
const input = ref<HTMLInputElement | HTMLTextAreaElement | null>(null)

const activate = () => {
  active.value = true
  text.value = model.value
  setTimeout(() => {
    if (input.value) {
      input.value.focus()
    }
    resize()
  }, 0)
}

const save = () => {
  model.value = text.value
  active.value = false
}

const resize = () => {
  if (input.value && props.type === 'textarea') {
    input.value.style.height = input.value.scrollHeight + 'px'
  }
}
</script>

<template>
  <div class="flex w-full items-center">
    <div
      v-if="!active"
      @click="activate"
      class="flex w-full cursor-pointer bg-transparent p-1 focus-visible:outline-none"
    >
      {{ model || placeholder }}
    </div>
    <div
      v-else
      class="flex w-full flex-row items-center border bg-transparent focus-visible:outline-none"
    >
      <div v-if="type === 'input'" class="flex w-full items-center">
        <input
          ref="input"
          autofocus
          v-model="text"
          :placeholder="placeholder"
          @keydown.enter="save"
          @blur="save"
          class="w-full border-none bg-transparent p-1 focus-visible:outline-none"
        />
      </div>
      <div v-else-if="type === 'textarea'" class="w-full">
        <textarea
          ref="input"
          v-model="text"
          :placeholder="placeholder"
          @keydown.enter="save"
          @blur="save"
          class="w-full border-none bg-transparent p-1 focus-visible:outline-none"
        />
      </div>
      <ShortCut class="mr-2 text-nowrap" keys="Press ↵ to save" />
    </div>
  </div>
</template>
