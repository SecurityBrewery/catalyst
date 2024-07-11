<script setup lang="ts">
import GrowTextarea from '@/components/form/GrowTextarea.vue'
import { Button } from '@/components/ui/button'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import { defineRule, useForm } from 'vee-validate'
import { watch } from 'vue'

import type { ReactionPython } from '@/lib/types'

const submitDisabled = defineModel('submitDisabled')

const props = defineProps<{
  reaction?: ReactionPython
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
    script: 'required'
  }
})

const equalReactionPython = (values: ReactionPython, reaction?: ReactionPython): boolean => {
  if (!reaction) return false

  return (
    reaction.name === values.name &&
    reaction.requirements === values.requirements &&
    reaction.script === values.script
  )
}

watch(
  values,
  () => {
    if (equalReactionPython(values, props.reaction)) {
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
    <FormField name="requirements" v-slot="{ componentField }">
      <FormItem>
        <FormLabel for="requirements" class="text-right">requirements.txt</FormLabel>
        <GrowTextarea id="requirements" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>
    <FormField name="script" v-slot="{ componentField }" validate-on-input>
      <FormItem>
        <FormLabel for="script" class="text-right">Script</FormLabel>
        <GrowTextarea id="script" class="col-span-3" v-bind="componentField" />
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
