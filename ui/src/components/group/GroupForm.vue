<script setup lang="ts">
import MultiSelect from '@/components/form/MultiSelect.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { useQuery } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { NewGroup } from '@/client'

const api = useAPI()

const submitDisabledReason = ref<string>('')

const props = defineProps<{
  group?: NewGroup
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

const { handleSubmit, validate, values } = useForm({
  initialValues: () => ({
    name: props.group?.name || '',
    permissions: props.group?.permissions || []
  }),
  validationSchema: {
    name: 'required'
  }
})

const equalGroup = (values: NewGroup, group?: NewGroup): boolean => {
  if (!group) return false

  return group.name === values.name && group.permissions === values.permissions
}

const updateSubmitDisabledReason = () => {
  if (isDemo.value) {
    submitDisabledReason.value = 'Groups cannot be created or edited in demo mode'

    return
  }

  if (equalGroup(values, props.group)) {
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
  () => props.group,
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
    name: values.name,
    permissions: values.permissions
  })
})

const permissionItems = computed(() => config.value?.permissions || [])
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

    <FormField
      key="permission.name"
      name="permissions"
      v-slot="{ componentField }"
      validate-on-input
    >
      <FormItem class="w-full">
        <div class="space-y-0.5">
          <FormLabel for="permissions" class="text-right">Permissions</FormLabel>
        </div>
        <FormControl>
          <MultiSelect
            v-bind="componentField"
            :items="permissionItems"
            placeholder="Select permissions..."
          />
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
              role="submit"
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
            <span v-else> Save the group. </span>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
      <slot name="cancel"></slot>
    </div>
  </form>
</template>
