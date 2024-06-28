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
import { toast } from '@/components/ui/toast'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Playbook } from '@/lib/types'

const queryClient = useQueryClient()
const router = useRouter()

const props = defineProps<{
  playbook: Playbook
}>()

const isOpen = ref(false)

const deletePlaybookMutation = useMutation({
  mutationFn: () => pb.collection('playbooks').delete(props.playbook.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['playbooks'] })
    router.push({ name: 'playbooks' })
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="outline">
        <Trash2 class="mr-2 h-4 w-4" />
        Delete Playbook
      </Button>
    </DialogTrigger>

    <DialogContent>
      <DialogHeader>
        <DialogTitle> Delete Playbook "{{ playbook.name }}"</DialogTitle>
        <DialogDescription> Are you sure you want to delete this playbook? </DialogDescription>
      </DialogHeader>

      <DialogFooter class="mt-2">
        <DialogClose as-child>
          <Button type="button" variant="secondary">Cancel</Button>
        </DialogClose>
        <Button type="button" variant="destructive" @click="deletePlaybookMutation.mutate">
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
