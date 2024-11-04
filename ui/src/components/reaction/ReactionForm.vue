<script setup lang="ts">
import ActionPythonFormFields from '@/components/reaction/ActionPythonFormFields.vue'
import ActionWebhookFormFields from '@/components/reaction/ActionWebhookFormFields.vue'
import TriggerHookFormFields from '@/components/reaction/TriggerHookFormFields.vue'
import TriggerScheduleFormFields from '@/components/reaction/TriggerScheduleFormFields.vue'
import TriggerWebhookFormFields from '@/components/reaction/TriggerWebhookFormFields.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { useQuery } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Reaction } from '@/lib/types'

const submitDisabledReason = ref<string>('')

const props = defineProps<{
  reaction?: Reaction
}>()

const emit = defineEmits(['submit'])

const isDemo = ref(false)

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: (): Promise<Record<string, Array<String>>> => pb.send('/api/config', {})
})

watch(
  () => config.value,
  () => {
    if (!config.value) return
    if (config.value['flags'].includes('demo')) {
      isDemo.value = true
    }
  },
  { immediate: true }
)

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

defineRule('triggerdata.expression', (value: string) => {
  if (values.trigger !== 'schedule') {
    return true
  }
  if (!value) {
    return 'This field is required'
  }
  const macros = ['@yearly', '@annually', '@monthly', '@weekly', '@daily', '@midnight', '@hourly']
  if (macros.includes(value)) {
    return true
  }
  const expression =
    /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]|1[0-9]|2[0-3])) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\*|([1-9]|1[0-2])|\*\/([1-9]|1[0-2])) (\*|([0-6])|\*\/([0-6]))$/
  if (value.match(expression)) {
    return true
  }
  return 'Invalid cron expression'
})

defineRule('triggerdata.token', (value: string) => {
  return true
})

defineRule('triggerdata.path', (value: string) => {
  if (values.trigger !== 'webhook') {
    return true
  }

  if (!value) {
    return 'This field is required'
  }

  const expression = /^[a-zA-Z0-9-_]+$/

  if (!value.match(expression)) {
    return 'Invalid path, only letters, numbers, dashes, and underscores are allowed'
  }

  return true
})

defineRule('triggerdata.collections', (value: string[]) => {
  if (values.trigger !== 'hook') {
    return true
  }

  if (!value) {
    return 'This field is required'
  }

  if (value.length === 0) {
    return 'At least one collection is required'
  }

  return true
})

defineRule('triggerdata.events', (value: string[]) => {
  if (values.trigger !== 'hook') {
    return true
  }

  if (!value) {
    return 'This field is required'
  }

  if (value.length === 0) {
    return 'At least one event is required'
  }

  return true
})

defineRule('actiondata.script', (value: string) => {
  if (values.action !== 'python') {
    return true
  }

  if (!value) {
    return 'This field is required'
  }

  return true
})

defineRule('actiondata.url', (value: string) => {
  if (values.action !== 'webhook') {
    return true
  }

  if (!value) {
    return 'This field is required'
  }

  if (!(value.startsWith('http://') || value.startsWith('https://'))) {
    return 'Invalid URL, must start with http:// or https://'
  }

  return true
})

const { handleSubmit, validate, values } = useForm({
  initialValues: props.reaction || {
    name: '',
    trigger: '',
    triggerdata: {},
    action: '',
    actiondata: {}
  },
  validationSchema: {
    name: 'required',
    trigger: 'required',
    'triggerdata.expression': 'triggerdata.expression',
    'triggerdata.token': 'triggerdata.token',
    'triggerdata.path': 'triggerdata.path',
    'triggerdata.collections': 'triggerdata.collections',
    'triggerdata.events': 'triggerdata.events',
    'actiondata.script': 'actiondata.script',
    'actiondata.url': 'actiondata.url',
    action: 'required'
  }
})

const equalReaction = (values: Reaction, reaction?: Reaction): boolean => {
  if (!reaction) return false

  return (
    reaction.name === values.name &&
    reaction.trigger === values.trigger &&
    JSON.stringify(reaction.triggerdata) === JSON.stringify(values.triggerdata) &&
    reaction.action === values.action &&
    JSON.stringify(reaction.actiondata) === JSON.stringify(values.actiondata)
  )
}

const updateSubmitDisabledReason = () => {
  if (isDemo.value) {
    submitDisabledReason.value = 'Reactions cannot be created or edited in demo mode'

    return
  }

  if (equalReaction(values, props.reaction)) {
    submitDisabledReason.value = 'Make changes to save'

    return
  }

  validate({ mode: 'silent' }).then((res) => {
    if (res.valid) {
      submitDisabledReason.value = ''
    } else {
      submitDisabledReason.value = 'Please fix the errors'
    }
  })
}

watch(
  () => isDemo.value,
  () => updateSubmitDisabledReason()
)

watch(
  () => props.reaction,
  () => updateSubmitDisabledReason(),
  { immediate: true }
)

watch(
  () => values,
  () => updateSubmitDisabledReason(),
  { deep: true, immediate: true }
)

const onSubmit = handleSubmit((values) => emit('submit', values))

const curlExample = computed(() => {
  let cmd = `curl`

  if (values.triggerdata.token) {
    cmd += ` -H "Authorization: Bearer ${values.triggerdata.token}"`
  }

  if (values.triggerdata.path) {
    cmd += ` https://${location.hostname}/reaction/${values.triggerdata.path}`
  }

  return cmd
})
</script>

<template>
  <form @submit="onSubmit" class="flex w-full flex-col items-start gap-4">
    <FormField name="name" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="name" class="text-right">Name</FormLabel>
        <Input id="name" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <Card class="w-full">
      <CardHeader>
        <CardTitle>Trigger</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <FormField name="trigger" v-slot="{ componentField }" validate-on-input>
          <FormItem>
            <FormLabel for="trigger" class="text-right">Type</FormLabel>
            <FormControl>
              <Select id="trigger" class="col-span-3" v-bind="componentField">
                <SelectTrigger class="font-medium">
                  <SelectValue placeholder="Select a type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="schedule">Schedule</SelectItem>
                    <SelectItem value="webhook">HTTP / Webhook</SelectItem>
                    <SelectItem value="hook">Collection Hook</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>
              <p>HTTP / Webhook: Receive a HTTP request.</p>
              <p>Collection Hook: Triggered by a collection and event.</p>
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <TriggerScheduleFormFields v-if="values.trigger === 'schedule'" />
        <TriggerWebhookFormFields v-else-if="values.trigger === 'webhook'" />
        <TriggerHookFormFields v-else-if="values.trigger === 'hook'" />

        <div v-if="values.trigger === 'webhook'">
          <Label for="url">Usage</Label>
          <Input id="url" readonly :modelValue="curlExample" class="bg-accent" />
        </div>
      </CardContent>
    </Card>

    <Card class="w-full">
      <CardHeader>
        <CardTitle>Action</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <FormField name="action" v-slot="{ componentField }" validate-on-input>
          <FormItem>
            <FormLabel for="action" class="text-right">Type</FormLabel>
            <FormControl>
              <Select id="action" class="col-span-3" v-bind="componentField">
                <SelectTrigger class="font-medium">
                  <SelectValue placeholder="Select a type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="python">Python</SelectItem>
                    <SelectItem value="webhook">HTTP / Webhook</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>
              <p>Python: Execute a Python script.</p>
              <p>HTTP / Webhook: Send an HTTP request.</p>
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <ActionPythonFormFields v-if="values.action === 'python'" />
        <ActionWebhookFormFields v-else-if="values.action === 'webhook'" />
      </CardContent>
    </Card>

    <Alert v-if="isDemo" variant="destructive">
      <AlertTitle>Cannot save</AlertTitle>
      <AlertDescription>{{ submitDisabledReason }}</AlertDescription>
    </Alert>
    <div class="flex gap-4">
      <TooltipProvider :delay-duration="0">
        <Tooltip>
          <TooltipTrigger class="cursor-default">
            <Button
              type="submit"
              :variant="submitDisabledReason !== '' ? 'secondary' : 'default'"
              :disabled="submitDisabledReason !== ''"
              :title="submitDisabledReason"
            >
              Save
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <span v-if="submitDisabledReason !== ''">
              {{ submitDisabledReason }}
            </span>
            <span v-else> Save the reaction. </span>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
      <slot name="cancel"></slot>
    </div>
  </form>
</template>
