<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DynamicInput from '@/components/input/DynamicInput.vue'
import PlaybookDeleteDialog from '@/components/playbook/PlaybookDeleteDialog.vue'
import StepView from '@/components/playbook/StepView.vue'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { toast } from '@/components/ui/toast'

import { Plus } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Playbook } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const {
  isPending,
  isError,
  data: playbook,
  error
} = useQuery({
  queryKey: ['playbooks', props.id],
  queryFn: (): Promise<Playbook> => pb.collection('playbooks').getOne(props.id)
})

const updatePlaybookMutation = useMutation({
  mutationFn: (update: any) => pb.collection('playbooks').update(props.id, update),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['playbooks'] }),
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

const updateName = (name: string) => {
  if (!playbook.value) return

  updatePlaybookMutation.mutate({ name: name })
}

const addStep = () => {
  if (!playbook.value) return

  const steps = JSON.parse(JSON.stringify(playbook.value.steps))
  steps.push({
    name: 'New Step',
    description: 'Description',
    type: 'task',
    schema: {}
  })
  updatePlaybookMutation.mutate({ steps })
}
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error" :value="playbook">
    <div class="flex h-full flex-1 flex-col overflow-hidden">
      <div class="flex items-center bg-background px-4 py-2">
        <div class="ml-auto">
          <PlaybookDeleteDialog v-if="playbook" :playbook="playbook" />
        </div>
      </div>
      <Separator />

      <ScrollArea v-if="playbook" class="flex-1">
        <div class="flex max-w-[640px] flex-col gap-4 p-4">
          <h1 class="text-3xl font-bold">
            <DynamicInput :modelValue="playbook.name" @update:modelValue="updateName" />
          </h1>
          <StepView
            v-for="(step, index) in playbook.steps"
            :key="index"
            :playbook="playbook"
            :step="step"
            :index="index"
          />

          <Button variant="outline" @click="addStep">
            <Plus class="mr-2 size-4" />
            Add Step
          </Button>
        </div>
      </ScrollArea>
    </div>
  </TanView>
</template>
