<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import UserSelect from '@/components/common/UserSelect.vue'
import DynamicInput from '@/components/input/DynamicInput.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TaskAddDialog from '@/components/ticket/task/TaskAddDialog.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'

import { Trash2, User2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Task, Ticket, User } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  tasks?: Array<Task>
}>()

const setTaskOwnerMutation = useMutation({
  mutationFn: (update: { id: string; user: User }): Promise<Task> =>
    pb.collection('tasks').update(update.id, {
      owner: update.user.id
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] }),
  onError: handleError
})

const update = (id: string, user: User) => setTaskOwnerMutation.mutate({ id, user })

const checkMutation = useMutation({
  mutationFn: (task: Task): Promise<Task> =>
    pb.collection('tasks').update(task.id, {
      open: !task.open
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] }),
  onError: handleError
})

const check = (task: Task) => checkMutation.mutate(task)

const updateTaskNameMutation = useMutation({
  mutationFn: (update: { id: string; name: string }): Promise<Task> =>
    pb.collection('tasks').update(update.id, {
      name: update.name
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] }),
  onError: handleError
})

const updateTaskName = (id: string, name: string) => updateTaskNameMutation.mutate({ id, name })
</script>

<template>
  <div class="mt-2 flex flex-col gap-2">
    <Card
      v-if="!tasks || tasks.length === 0"
      class="flex h-10 items-center p-4 text-muted-foreground"
    >
      No tasks added yet.
    </Card>
    <Card v-else>
      <PanelListElement v-for="task in tasks" :key="task.id" class="pr-1">
        <div class="flex flex-row items-center">
          <Checkbox :checked="!task.open" class="mr-2" @click="check(task)" />
          <DynamicInput
            :modelValue="task.name"
            @update:modelValue="updateTaskName(task.id, $event)"
            class="mr-2 flex-1"
          />
        </div>
        <div class="ml-auto flex items-center">
          <UserSelect v-if="!task.expand.owner" @update:modelValue="update(task.id, $event)">
            <Button variant="outline" role="combobox" class="h-8">
              <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
              Unassigned
            </Button>
          </UserSelect>
          <UserSelect
            v-else
            :modelValue="task.expand.owner"
            @update:modelValue="update(task.id, $event)"
          >
            <Button variant="outline" role="combobox" class="mr-2 h-8">
              <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
              {{ task.expand.owner.name }}
            </Button>
          </UserSelect>
          <DeleteDialog
            v-if="task"
            collection="tasks"
            :id="task.id"
            :name="task.name"
            :singular="'Task'"
            :queryKey="['tickets', ticket.id]"
          >
            <Button variant="ghost" size="icon" class="h-8 w-8">
              <Trash2 class="size-4" />
            </Button>
          </DeleteDialog>
        </div>
      </PanelListElement>
    </Card>
    <TaskAddDialog :ticket="ticket" class="w-full" />
  </div>
</template>
