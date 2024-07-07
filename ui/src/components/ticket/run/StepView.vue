<script setup lang="ts">
import JSONSchemaFormFields from '@/components/form/JSONSchemaFormFields.vue'
import MarkdownView from '@/components/input/MarkdownView.vue'
import StatusIcon from '@/components/ticket/StatusIcon.vue'
import { Button } from '@/components/ui/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { Separator } from '@/components/ui/separator'
import { toast } from '@/components/ui/toast'

import { ChevronDown } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { computed, ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Run, RunStep } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  run: Run
  index: number
  step: RunStep
}>()

const isOpen = ref(false)
// const data = ref({})
const state = ref(props.step.state)
const error = ref<string | undefined>(undefined)

const updateRunMutation = useMutation({
  mutationFn: (update: any) => pb.collection('runs').update(props.run.id, update),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets', props.run.ticket] }),
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const validationSchema = computed(() => {
  const fields: Record<string, any> = {
    name: 'required'
  }

  if (props.step.schema.properties) {
    Object.keys(props.step.schema.properties).forEach((key) => {
      const property = props.step.schema.properties[key]
      if (property.type === 'string') {
        fields[key] =
          props.step.schema.required && props.step.schema.required.includes(key) ? 'required' : ''
      } else if (property.type === 'boolean') {
        fields[key] = ''
      }
    })
  }

  return fields
})

const { validate } = useForm({
  validationSchema: validationSchema.value
})

const toggleStatus = () => {
  if (props.step.type !== 'jsonschema') {
    updateState()

    return
  }

  validate().then((res) => {
    if (!res.valid) {
      let err = ''
      for (const [key, value] of Object.entries(res.errors)) {
        err += `${key}: ${value}\n`
      }
      error.value = err

      return
    }

    updateState()
  })
}

const updateState = () => {
  const steps = JSON.parse(JSON.stringify(props.run.steps))
  if (props.step.status === 'open') {
    steps[props.index].status = 'completed'
    steps[props.index].state = state.value
  } else {
    steps[props.index].status = 'open'
  }
  updateRunMutation.mutate({ steps })
}
</script>

<template>
  <Collapsible v-model:open="isOpen" class="w-full space-y-2">
    <div class="flex items-center justify-between">
      <CollapsibleTrigger as-child>
        <div
          class="flex w-full cursor-pointer items-center justify-between border-t px-4 py-2 transition-transform duration-200"
          :class="{ 'border-b bg-accent': isOpen }"
        >
          <div class="flex items-center gap-2 space-x-2">
            <ChevronDown
              class="h-4 w-4 transition-transform duration-200"
              :class="{ 'rotate-180': isOpen }"
            />
            {{ step.name }}
          </div>
          <div class="flex items-center gap-2 space-x-2">
            <span class="sr-only">Toggle</span>

            <StatusIcon :status="step.status" class="h-5 w-5" />
          </div>
        </div>
      </CollapsibleTrigger>
    </div>

    <CollapsibleContent class="space-y-2 px-2 pb-2">
      <MarkdownView :markdown="step.description" />

      <Separator v-if="step.type && step.type === 'jsonschema' && step.schema" />

      <JSONSchemaFormFields v-model="state" :schema="step.schema" />

      <Button variant="default" size="sm" class="w-full" @click="toggleStatus">
        {{ step.status === 'completed' ? 'Reopen' : 'Close' }} Step
      </Button>
    </CollapsibleContent>
  </Collapsible>
</template>
