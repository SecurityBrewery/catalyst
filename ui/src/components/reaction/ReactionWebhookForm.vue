<script setup lang="ts">
import GrowTextarea from '@/components/form/GrowTextarea.vue'
import { Button } from '@/components/ui/button'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import { defineRule, useForm } from 'vee-validate'
import { watch } from 'vue'

import type { ReactionWebhook } from '@/lib/types'

const submitDisabled = defineModel('submitDisabled')

const props = defineProps<{
  reaction?: ReactionWebhook
}>()

const emit = defineEmits(['submit'])

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const { handleSubmit, validate, values } = useForm({
  initialValues: props.reaction,
  validationSchema: {
    name: 'required',
    destination: 'required'
  }
})

const equalReactionWebhook = (values: ReactionWebhook, reaction?: ReactionWebhook): boolean => {
  if (!reaction) return false

  return (
    reaction.name === values.name &&
    reaction.headers === values.headers &&
    reaction.destination === values.destination
  )
}

watch(
  values,
  () => {
    if (equalReactionWebhook(values, props.reaction)) {
      submitDisabled.value = true

      return
    }

    validate({ mode: 'silent' }).then((res) => (submitDisabled.value = !res.valid))
  },
  { deep: true, immediate: true }
)

const onSubmit = handleSubmit((values) => emit('submit', values))
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
    <FormField name="headers" v-slot="{ componentField }">
      <FormItem>
        <FormLabel for="headers" class="text-right">Headers (one header per line)</FormLabel>
        <GrowTextarea
          id="headers"
          class="col-span-3"
          v-bind="componentField"
          placeholder="Content-Type: application/json"
        />
        <FormMessage />
      </FormItem>
    </FormField>
    <FormField name="destination" v-slot="{ componentField }" validate-on-input>
      <FormItem>
        <FormLabel for="destination" class="text-right">Destination</FormLabel>
        <Input
          id="destination"
          class="col-span-3"
          v-bind="componentField"
          placeholder="https://example.com/webhook"
        />
        <FormMessage />
      </FormItem>
    </FormField>
    <slot>
      <Button
        type="submit"
        :variant="submitDisabled ? 'secondary' : 'default'"
        :disabled="submitDisabled"
        >Save</Button
      >
    </slot>
  </form>
</template>
