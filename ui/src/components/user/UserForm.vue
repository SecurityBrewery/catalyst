<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { useQuery } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { UserUpdate } from '@/client/models'

const api = useAPI()

const submitDisabledReason = ref<string>('')

const props = defineProps<{
  user?: UserUpdate
}>()

const emit = defineEmits(['submit'])

const isDemo = ref(false)

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: () => api.getConfig()
})

watch(
  () => config.value,
  () => {
    if (!config.value) return
    if (config.value.flags.includes('demo')) {
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

defineRule('username', (value: string) => {
  const usernamePattern = /^[a-z0-9_]{3,20}$/
  if (!usernamePattern.test(value)) {
    return 'Username must be 3-20 characters long and can only contain lowercase letters, numbers, and underscores'
  }

  return true
})

defineRule('email', (value: string) => {
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailPattern.test(value)) {
    return 'Please enter a valid email address'
  }

  return true
})

const { handleSubmit, validate, values } = useForm({
  initialValues: props.user || {
    username: '',
    avatar: '',
    email: '',
    emailVisibility: true,
    name: '',
    verified: false
  },
  validationSchema: {
    username: 'username',
    email: 'email',
    name: 'required'
  }
})

const equalUser = (values: UserUpdate, user?: UserUpdate): boolean => {
  if (!user) return false

  return (
    user.username === values.username &&
    user.avatar === values.avatar &&
    user.email === values.email &&
    user.emailVisibility === values.emailVisibility &&
    user.name === values.name &&
    user.verified === values.verified
  )
}

const updateSubmitDisabledReason = () => {
  if (props.user?.username === 'system') {
    submitDisabledReason.value = 'The system user cannot be modified'

    return
  }

  if (isDemo.value) {
    submitDisabledReason.value = 'Users cannot be created or edited in demo mode'

    return
  }

  if (equalUser(values, props.user)) {
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
  () => props.user,
  () => updateSubmitDisabledReason(),
  { immediate: true }
)

watch(
  () => values,
  () => updateSubmitDisabledReason(),
  { deep: true, immediate: true }
)

const onSubmit = handleSubmit((values) => emit('submit', values))
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

    <FormField name="username" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="username" class="text-right">Username</FormLabel>
        <Input id="username" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="email" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="email" class="text-right">Email</FormLabel>
        <Input id="email" type="email" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="verified" v-slot="{ value, handleChange }">
      <FormItem class="w-full items-center gap-2">
        <FormLabel>Verified</FormLabel>
        <div class="flex flex-row items-center gap-2">
          <FormControl>
            <Switch :checked="value" @update:checked="handleChange" />
          </FormControl>
          <FormDescription> Check to allow the user to log in. </FormDescription>
        </div>
        <FormMessage />
      </FormItem>
    </FormField>

    <Alert v-if="props.user?.username === 'system'" variant="destructive">
      <AlertTitle>Cannot save</AlertTitle>
      <AlertDescription>The system user cannot be modified.</AlertDescription>
    </Alert>
    <Alert v-else-if="isDemo" variant="destructive">
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
            <span v-else> Save the user. </span>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
      <slot name="cancel"></slot>
    </div>
  </form>
</template>
