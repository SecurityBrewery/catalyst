<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

defineProps<{
  modelValue: boolean
  imageUrl: string
  fileName: string
}>()

const emit = defineEmits(['update:modelValue'])

const closeModal = () => {
  emit('update:modelValue', false)
}
</script>

<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    @click.self="closeModal"
  >
    <div class="relative rounded-lg bg-white p-4 shadow-xl">
      <img :src="imageUrl" :alt="fileName" class="max-h-[80vh] max-w-[80vw]" />
      <Button
        variant="ghost"
        size="icon"
        class="absolute -top-1 -right-1 h-8 w-8 rounded-full bg-white"
        @click="closeModal"
      >
        <X class="h-5 w-5" />
        <span class="sr-only">Close image modal</span>
      </Button>
    </div>
  </div>
</template>
