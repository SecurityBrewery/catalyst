<script setup lang="ts">
import DynamicInput from '@/components/input/DynamicInput.vue'
import DynamicMDEditor from '@/components/input/DynamicMDEditor.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
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
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { toast } from '@/components/ui/toast'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Playbook, PlaybookStep } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  index: number
  playbook: Playbook
  step: PlaybookStep
}>()

const isOpen = ref(false)
const message = ref(props.step.description)

const updatePlaybookMutation = useMutation({
  mutationFn: (update: any) => pb.collection('playbooks').update(props.playbook.id, update),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['playbooks'] })
    isOpen.value = false
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

const updateStepName = (index: number, name: string) => {
  const steps = JSON.parse(JSON.stringify(props.playbook.steps))
  steps[index].name = name
  updatePlaybookMutation.mutate({ steps })
}

const updateStepType = (index: number, type: string) => {
  const steps = JSON.parse(JSON.stringify(props.playbook.steps))
  steps[index].type = type
  updatePlaybookMutation.mutate({ steps })
}

const deleteStep = (index: number) => {
  const steps = JSON.parse(JSON.stringify(props.playbook.steps))
  steps.splice(index, 1)
  updatePlaybookMutation.mutate({ steps })
}

const saveDescription = () => {
  const steps = JSON.parse(JSON.stringify(props.playbook.steps))
  steps[props.index].description = message.value
  updatePlaybookMutation.mutate({ steps })
}

const updateStepSchema = (index: number, schema: string) => {
  const steps = JSON.parse(JSON.stringify(props.playbook.steps))
  steps[index].schema = JSON.parse(schema)
  updatePlaybookMutation.mutate({ steps })
}

const editDescription = ref(true)
</script>

<template>
  <Card>
    <CardHeader class="relative">
      <CardTitle class="flex flex-row items-center">
        <DynamicInput :modelValue="step.name" @update:modelValue="updateStepName(index, $event)" />

        <!-- Select :modelValue="step.type" @update:modelValue="updateStepType(index, $event)">
          <SelectTrigger class="mr-2 w-40">
            <SelectValue placeholder="Select Step Type" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Step Type</SelectLabel>
              <SelectItem value="task"> Task</SelectItem>
              <SelectItem value="jsonschema"> JSONSchema</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select -->
        <Dialog v-model:open="isOpen">
          <DialogTrigger as-child>
            <Button variant="outline">
              <Trash2 class="mr-2 h-4 w-4" />
              Delete Step
            </Button>
          </DialogTrigger>

          <DialogContent>
            <DialogHeader>
              <DialogTitle> Delete Step "{{ step.name }}"</DialogTitle>
              <DialogDescription> Are you sure you want to delete this step?</DialogDescription>
            </DialogHeader>

            <DialogFooter class="mt-2">
              <DialogClose as-child>
                <Button type="button" variant="secondary">Cancel</Button>
              </DialogClose>
              <Button type="button" variant="destructive" @click="deleteStep(index)">
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </CardTitle>
      <CardDescription>
        <div class="rounded border p-4">
          <DynamicMDEditor
            v-model="message"
            v-model:edit="editDescription"
            hide-cancel
            @save="saveDescription"
          />
          <!-- MDEditor v-model="message" @save="saveDescription" hide-cancel / -->
        </div>
      </CardDescription>
      <div class="flex flex-row items-center gap-2"></div>
    </CardHeader>

    <!-- CardContent v-if="step.type === 'jsonschema'" class="relative">
      <Textarea
        :modelValue="JSON.stringify(step.schema, null, 2)"
        @update:modelValue="updateStepSchema(index, $event as any)"
        class="w-full"
        rows="10"
        placeholder="Enter JSON Schema"
      />
    </CardContent -->
  </Card>
</template>
