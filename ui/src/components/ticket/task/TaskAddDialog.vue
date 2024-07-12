<script setup lang="ts">
import ShortCut from '@/components/ShortCut.vue'
import UserSelect from '@/components/common/UserSelect.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { Plus, User2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Task, Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const name = ref('')
const owner = ref(pb.authStore.model)
const isOpen = ref(false)

const addTaskMutation = useMutation({
  mutationFn: (): Promise<Task> => {
    if (!pb.authStore.model) return Promise.reject('Not authenticated')
    return pb.collection('tasks').create({
      ticket: props.ticket.id,
      name: name.value,
      open: true,
      owner: owner.value ? owner.value.id : pb.authStore.model.id
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    name.value = ''
    owner.value = pb.authStore.model
    isOpen.value = false
  },
  onError: handleError
})

const submitDisabled = computed(() => !name.value || !owner.value)
</script>

<template>
  <Button v-if="!isOpen" variant="outline" @click="isOpen = true">
    <Plus class="mr-2 size-4" />
    Add Task
  </Button>
  <Card v-else class="flex flex-row items-center gap-2 p-2">
    <Input
      v-model="name"
      name="name"
      placeholder="Add a task..."
      autofocus
      @keydown.meta.enter="addTaskMutation.mutate"
      @keydown.ctrl.enter="addTaskMutation.mutate"
    />
    <UserSelect v-model="owner as any">
      <Button variant="outline" role="combobox">
        <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
        {{ owner?.name }}
      </Button>
    </UserSelect>
    <Button variant="outline" @click="isOpen = false"> Cancel</Button>
    <Button :disabled="submitDisabled" @click="addTaskMutation.mutate">
      Save
      <ShortCut keys="⌘ ↵" />
    </Button>
  </Card>
</template>
