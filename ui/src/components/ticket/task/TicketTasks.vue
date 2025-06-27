<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import UserSelect from '@/components/common/UserSelect.vue'
import DynamicInput from '@/components/input/DynamicInput.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TaskAddDialog from '@/components/ticket/task/TaskAddDialog.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { useToast } from '@/components/ui/toast/use-toast'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'

import { useAPI } from '@/api'
import type { ExtendedTask, Task, Ticket, User } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  ticket: Ticket
  tasks?: Array<ExtendedTask>
}>()

const setTaskOwnerMutation = useMutation({
  mutationFn: (update: { id: string; owner: string }): Promise<Task> =>
    api.updateTask({
      id: update.id,
      taskUpdate: {
        owner: update.owner
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tasks', props.ticket.id] })
    toast({
      title: 'Owner updated',
      description: 'The task owner has been updated'
    })
  },
  onError: handleError('Failed to update task owner')
})

const update = (id: string, user: User) => setTaskOwnerMutation.mutate({ id, owner: user.id })

const checkMutation = useMutation({
  mutationFn: (task: Task): Promise<Task> =>
    api.updateTask({
      id: task.id,
      taskUpdate: {
        open: !task.open
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tasks', props.ticket.id] })
    toast({
      title: 'Task updated',
      description: 'The task status has been updated'
    })
  },
  onError: handleError('Failed to update task status')
})

const check = (task: Task) => checkMutation.mutate(task)

const updateTaskNameMutation = useMutation({
  mutationFn: (update: { id: string; name: string }): Promise<Task> =>
    api.updateTask({
      id: update.id,
      taskUpdate: {
        name: update.name
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tasks', props.ticket.id] })
    toast({
      title: 'Task updated',
      description: 'The task name has been updated'
    })
  },
  onError: handleError('Failed to update task name')
})

const deleteMutation = useMutation({
  mutationFn: (id: string) => api.deleteTask({ id }),
  onSuccess: (data, id) => {
    queryClient.removeQueries({ queryKey: ['tasks', id] })
    queryClient.invalidateQueries({ queryKey: ['tasks', props.ticket.id] })
    toast({
      title: 'Task deleted',
      description: 'The task has been deleted successfully'
    })
  },
  onError: handleError('Failed to delete task')
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
        <div class="ml-auto flex items-center gap-1">
          <UserSelect
            :userID="task.owner"
            :userName="task.ownerName"
            @select="update(task.id, $event)"
          />
          <DeleteDialog
            v-if="task"
            :name="task.name"
            singular="Task"
            @delete="deleteMutation.mutate(task.id)"
          >
            <Button variant="ghost" size="icon" class="size-10">
              <Trash2 class="size-4" />
              <span class="sr-only">Delete Task</span>
            </Button>
          </DeleteDialog>
        </div>
      </PanelListElement>
    </Card>
    <TaskAddDialog :ticket="ticket" class="w-full" />
  </div>
</template>
