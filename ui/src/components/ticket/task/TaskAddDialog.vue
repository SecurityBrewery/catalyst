<script setup lang="ts">
import ShortCut from '@/components/ShortCut.vue'
import UserSelect from '@/components/common/UserSelect.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useToast } from '@/components/ui/toast/use-toast'

import { Plus, User2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

import { useAPI } from '@/api'
import type { Task, Ticket, User } from '@/client/models'
import { handleError } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const api = useAPI()

const queryClient = useQueryClient()
const authStore = useAuthStore()
const { toast } = useToast()

const props = defineProps<{
  ticket: Ticket
}>()

const name = ref('')
const ownerID = ref<string | undefined>(authStore.user?.id)
const ownerName = ref<string | undefined>(authStore.user?.name)
const isOpen = ref(false)

const addTaskMutation = useMutation({
  mutationFn: (): Promise<Task> => {
    if (!authStore.user) return Promise.reject('Not authenticated')
    return api.createTask({
      newTask: {
        ticket: props.ticket.id,
        name: name.value,
        open: true,
        owner: ownerID.value ? ownerID.value : ''
      }
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tasks', props.ticket.id] })
    toast({
      title: 'Task created',
      description: 'The task has been created successfully'
    })
    name.value = ''
    isOpen.value = false
  },
  onError: handleError
})

const submitDisabled = computed(() => !name.value || !ownerID.value)

const select = (user: User) => {
  ownerID.value = user.id
  ownerName.value = user.name
}
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
    <UserSelect :userID="ownerID" :userName="ownerName" @select="select" />
    <Button variant="outline" @click="isOpen = false"> Cancel</Button>
    <Button :disabled="submitDisabled" @click="addTaskMutation.mutate">
      Save
      <ShortCut keys="⌘ ↵" />
    </Button>
  </Card>
</template>
