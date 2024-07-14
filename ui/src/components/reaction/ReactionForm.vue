<script setup lang="ts">
import ReactionPythonFormFields from '@/components/reaction/ReactionPythonFormFields.vue'
import ReactionWebhookFormFields from '@/components/reaction/ReactionWebhookFormFields.vue'
import TriggerHookFormFields from '@/components/reaction/TriggerHookFormFields.vue'
import TriggerWebhookFormFields from '@/components/reaction/TriggerWebhookFormFields.vue'
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

import { defineRule, useForm } from 'vee-validate'
import { computed, watch } from 'vue'

import type { Reaction } from '@/lib/types'

const submitDisabled = defineModel('submitDisabled')

const props = defineProps<{
  reaction?: Reaction
}>()

const emit = defineEmits(['submit'])

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
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

defineRule('reactiondata.script', (value: string) => {
  if (values.reaction !== 'python') {
    return true
  }

  if (!value) {
    return 'This field is required'
  }

  return true
})

defineRule('reactiondata.url', (value: string) => {
  if (values.reaction !== 'webhook') {
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
    reaction: '',
    reactiondata: {}
  },
  validationSchema: {
    name: 'required',
    trigger: 'required',
    'triggerdata.token': 'triggerdata.token',
    'triggerdata.path': 'triggerdata.path',
    'triggerdata.collections': 'triggerdata.collections',
    'triggerdata.events': 'triggerdata.events',
    'reactiondata.script': 'reactiondata.script',
    'reactiondata.url': 'reactiondata.url',
    reaction: 'required'
  }
})

const equalReaction = (values: Reaction, reaction?: Reaction): boolean => {
  if (!reaction) return false

  return (
    reaction.name === values.name &&
    reaction.trigger === values.trigger &&
    JSON.stringify(reaction.triggerdata) === JSON.stringify(values.triggerdata) &&
    reaction.reaction === values.reaction &&
    JSON.stringify(reaction.reactiondata) === JSON.stringify(values.reactiondata)
  )
}

watch(
  () => props.reaction,
  () => {
    if (equalReaction(values, props.reaction)) {
      submitDisabled.value = true
    }
  },
  { immediate: true }
)

watch(
  values,
  () => {
    if (equalReaction(values, props.reaction)) {
      submitDisabled.value = true

      return
    }

    validate({ mode: 'silent' }).then((res) => {
      submitDisabled.value = !res.valid
    })
  },
  { deep: true, immediate: true }
)

const onSubmit = handleSubmit((values) => emit('submit', values))

const curlExample = computed(() => {
  let cmd = `curl`

  if (values.triggerdata.token) {
    cmd += ` -H "Auth: Bearer ${values.triggerdata.token}"`
  }

  if (values.triggerdata.path) {
    cmd += ` https://${location.hostname}/reaction/${values.triggerdata.path}`
  }

  return cmd
})
</script>

<template>
  <form @submit="onSubmit" class="flex flex-col gap-4">
    <FormField name="name" v-slot="{ componentField }" validate-on-input>
      <FormItem>
        <FormLabel for="name" class="text-right">Name</FormLabel>
        <Input id="name" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <Card>
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

        <TriggerWebhookFormFields v-if="values.trigger === 'webhook'" />
        <TriggerHookFormFields v-else-if="values.trigger === 'hook'" />

        <div v-if="values.trigger === 'webhook'">
          <Label for="url">Usage</Label>
          <Input id="url" readonly :modelValue="curlExample" class="bg-accent" />
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>Reaction</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <FormField name="reaction" v-slot="{ componentField }" validate-on-input>
          <FormItem>
            <FormLabel for="reaction" class="text-right">Type</FormLabel>
            <FormControl>
              <Select id="reaction" class="col-span-3" v-bind="componentField">
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

        <ReactionPythonFormFields v-if="values.reaction === 'python'" />
        <ReactionWebhookFormFields v-else-if="values.reaction === 'webhook'" />
      </CardContent>
    </Card>

    <slot>
      <Button
        type="submit"
        :variant="submitDisabled ? 'secondary' : 'default'"
        :disabled="submitDisabled"
        >Save
      </Button>
    </slot>
  </form>
</template>
