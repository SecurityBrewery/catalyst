<script setup lang="ts">
import GrowTextarea from '@/components/form/GrowTextarea.vue'
import IconPicker from '@/components/form/IconPicker.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { useQuery } from '@tanstack/vue-query'
import isEqual from 'lodash.isequal'
import { defineRule, useForm } from 'vee-validate'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { NewType } from '@/client/models'

const api = useAPI()

const submitDisabledReason = ref<string>('')

const props = defineProps<{
  type?: NewType
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

defineRule('validate_json', (value: string) => {
  try {
    JSON.parse(value)
    return true
  } catch (e) {
    return 'Invalid JSON format'
  }
})

const { handleSubmit, validate, values } = useForm({
  initialValues: () => ({
    singular: props.type?.singular || '',
    plural: props.type?.plural || '',
    icon: props.type?.icon || '',
    schema: JSON.stringify(props.type?.schema, null, 2) || '{}'
  }),
  validationSchema: {
    singular: 'required',
    plural: 'required',
    schema: 'validate_json'
  }
})

const equalType = (values: NewType, type?: NewType): boolean => {
  if (!type) return false

  return (
    type.singular === values.singular &&
    type.plural === values.plural &&
    type.icon === values.icon &&
    isEqual(type.schema, JSON.parse(values.schema))
  )
}

const updateSubmitDisabledReason = () => {
  if (isDemo.value) {
    submitDisabledReason.value = 'Types cannot be created or edited in demo mode'

    return
  }

  if (equalType(values, props.type)) {
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
  () => props.type,
  () => updateSubmitDisabledReason(),
  { immediate: true }
)

watch(
  () => values,
  () => updateSubmitDisabledReason(),
  { deep: true, immediate: true }
)

const onSubmit = handleSubmit((values) => {
  emit('submit', {
    singular: values.singular,
    plural: values.plural,
    icon: values.icon,
    schema: JSON.parse(values.schema)
  })
})
</script>

<template>
  <form @submit="onSubmit" class="flex w-full flex-col items-start gap-4">
    <FormField name="singular" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="singular" class="text-left">Singular</FormLabel>
        <Input id="singular" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="plural" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="plural" class="text-left">Plural</FormLabel>
        <Input id="plural" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="icon" v-slot="{ componentField }" validate-on-input>
      <FormItem class="flex w-full flex-col">
        <FormLabel for="icon" class="text-start">Icon</FormLabel>
        <span class="text-xs text-muted-foreground">
          Select a suggested icon or browse
          <a class="text-blue-500" href="https://lucide.dev/icons/" target="_blank">lucide.dev</a>.
        </span>
        <IconPicker id="icon" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="schema" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="schema" class="text-start">Schema</FormLabel>
        <FormControl>
          <GrowTextarea id="schema" class="col-span-3" v-bind="componentField" />
        </FormControl>
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
              Save
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <span v-if="submitDisabledReason !== ''">
              {{ submitDisabledReason }}
            </span>
            <span v-else> Save the type. </span>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
      <slot name="cancel"></slot>
    </div>
  </form>
</template>
