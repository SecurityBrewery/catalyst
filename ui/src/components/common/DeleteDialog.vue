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

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { type RouteLocationRaw, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()
const router = useRouter()

const props = defineProps<{
  collection: string
  id: string
  name: string
  singular: string
  queryKey: string[]
  to?: RouteLocationRaw
}>()

const isOpen = ref(false)

const deleteMutation = useMutation({
  mutationFn: () => pb.collection(props.collection).delete(props.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: props.queryKey })
    if (props.to) router.push(props.to)
  },
  onError: handleError
})
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
        <Button type="button" variant="destructive" @click="deleteMutation.mutate"> Delete </Button>
        <DialogClose as-child>
          <Button type="button" variant="secondary">Cancel</Button>
        </DialogClose>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
