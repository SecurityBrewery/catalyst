<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { useQuery } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'

const api = useAPI()

const submitDisabledReason = ref<string>('')

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

defineRule('password', (value: string) => {
  if (!value) {
    return 'Password is required'
  }

  if (value.length < 8) {
    return 'Password must be at least 8 characters long'
  }

  return true
})

defineRule('passwordConfirm', (value: string) => {
  if (!value) {
    return 'Password confirmation is required'
  }

  if (value !== values.password) {
    return 'Passwords do not match'
  }

  return true
})

const { handleSubmit, validate, values, resetForm } = useForm({
  initialValues: {
    password: '',
    passwordConfirm: ''
  },
  validationSchema: {
    password: 'password',
    passwordConfirm: 'passwordConfirm'
  }
})

const reset = () =>
  resetForm({
    values: { password: '', passwordConfirm: '' }
  })

defineExpose({ reset })

const updateSubmitDisabledReason = () => {
  if (isDemo.value) {
    submitDisabledReason.value = 'Users cannot be created or edited in demo mode'

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
  () => values,
  () => updateSubmitDisabledReason(),
  { deep: true, immediate: true }
)

const onSubmit = handleSubmit((values) => emit('submit', values))
</script>

<template>
  <form @submit="onSubmit" class="flex w-full flex-col items-start gap-4">
    <FormField name="password" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="password" class="text-right">Password</FormLabel>
        <Input id="password" type="password" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="passwordConfirm" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="passwordConfirm" class="text-right">Confirm Password</FormLabel>
        <Input id="passwordConfirm" type="password" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

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
              Set Password
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
