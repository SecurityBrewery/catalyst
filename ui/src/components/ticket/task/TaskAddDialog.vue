<script setup lang="ts">
import ShortCut from '@/components/ShortCut.vue'
import UserSelect from '@/components/common/UserSelect.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { Plus, User2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

import { api } from '@/api'
import type { Task, Ticket, User } from '@/client/models'
import { handleError } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const queryClient = useQueryClient()
const authStore = useAuthStore()

const props = defineProps<{
  ticket: Ticket
}>()

const name = ref('')
const owner = ref<string | undefined>(authStore.user?.id)
const isOpen = ref(false)

const addTaskMutation = useMutation({
  mutationFn: (): Promise<Task> => {
    if (!authStore.user) return Promise.reject('Not authenticated')
    return api.createTask({
      newTask: {
        ticket: props.ticket.id,
        name: name.value,
        open: true,
        owner: owner.value ? owner.value : authStore.user?.id
      }
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    name.value = ''
    owner.value = authStore.user?.id
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
    <UserSelect v-model="owner">
      <Button variant="outline" role="combobox">
        <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
        {{ owner }}
        <!-- TODO -->
      </Button>
    </UserSelect>
    <Button variant="outline" @click="isOpen = false"> Cancel</Button>
    <Button :disabled="submitDisabled" @click="addTaskMutation.mutate">
      Save
      <ShortCut keys="⌘ ↵" />
    </Button>
  </Card>
</template>
