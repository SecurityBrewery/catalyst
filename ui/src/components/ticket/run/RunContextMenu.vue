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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { toast } from '@/components/ui/toast'

import { MoreVertical, Trash } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Run } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  run: Run
}>()

const isOpen = ref(false)

const removeRunMutation = useMutation({
  mutationFn: (): Promise<boolean> => pb.collection('runs').delete(props.run.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.run.ticket] })
    isOpen.value = false
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
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button size="icon" variant="outline" class="h-8 w-8">
          <MoreVertical class="h-3.5 w-3.5" />
          <span class="sr-only">More</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem class="cursor-pointer" as-child>
          <DialogTrigger class="w-full">
            <Trash class="mr-2 h-4 w-4" />
            Delete
          </DialogTrigger>
        </DropdownMenuItem>
      </DropdownMenuContent>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete run</DialogTitle>
          <DialogDescription> Are you sure you want to delete this run?</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose as-child>
            <Button type="button" variant="secondary"> Cancel</Button>
          </DialogClose>
          <Button @click="removeRunMutation.mutate" variant="destructive"> Delete </Button>
        </DialogFooter>
      </DialogContent>
    </DropdownMenu>
  </Dialog>
</template>
