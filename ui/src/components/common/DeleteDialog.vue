<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'

import { Trash2 } from 'lucide-vue-next'

import { ref } from 'vue'

const props = defineProps<{
  name: string
  singular: string
}>()

const isOpen = ref(false)

const emit = defineEmits<{
  (e: 'delete'): void
}>()

const deleteRecord = () => {
  emit('delete')
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <slot>
        <Button variant="outline">
          <Trash2 class="mr-2 h-4 w-4" />
          Delete {{ props.singular }}
        </Button>
      </slot>
    </DialogTrigger>

    <DialogContent>
      <DialogHeader>
        <DialogTitle> Delete {{ props.singular }} "{{ props.name }}"</DialogTitle>
        <DialogDescription>
          Are you sure you want to delete this {{ props.singular }}?</DialogDescription
        >
      </DialogHeader>

      <DialogFooter class="mt-2 sm:justify-start">
        <Button type="button" variant="destructive" @click="deleteRecord"> Delete </Button>
        <DialogClose as-child>
          <Button type="button" variant="secondary">Cancel</Button>
        </DialogClose>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
