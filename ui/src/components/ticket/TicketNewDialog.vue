<script setup lang="ts">
import JSONSchemaFormFields from '@/components/form/JSONSchemaFormFields.vue'
import { Button } from '@/components/ui/button'
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
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket, Type } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()
const router = useRouter()

const props = defineProps<{
  selectedType: Type
}>()

const isOpen = ref(false)

const addTicketMutation = useMutation({
  mutationFn: (): Promise<Ticket> => {
    if (!pb.authStore.model) return Promise.reject('Not logged in')

    return pb.collection('tickets').create({
      name: name.value,
      description: description.value,
      open: true,
      type: props.selectedType.id,
      schema: props.selectedType.schema,
      state: state.value,
      owner: pb.authStore.model.id
    })
  },
  onSuccess: (data: Ticket) => {
    router.push(`/tickets/${props.selectedType.id}/${data.id}`)
    queryClient.invalidateQueries({ queryKey: ['tickets'] })
    isOpen.value = false
  },
  onError: handleError
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

  Object.keys(props.selectedType.schema.properties).forEach((key) => {
    const property = props.selectedType.schema.properties[key]
    if (property.type === 'string') {
      fields[key] =
        props.selectedType.schema.required && props.selectedType.schema.required.includes(key)
          ? 'required'
          : ''
    } else if (property.type === 'boolean') {
      fields[key] = ''
    }
  })

  return fields
})

const { handleSubmit, validate } = useForm({
  validationSchema: validationSchema.value
})

const onSubmit = handleSubmit((values: any) => {
  validate().then((res) => {
    if (res.valid) {
      addTicketMutation.mutate(values)
    }
  })
})

const state = ref({})
const name = ref('')
const description = ref('')

watch(
  () => isOpen.value,
  () => {
    if (isOpen.value) {
      name.value = ''
      description.value = ''
      state.value = {}
    }
  }
)
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="ghost"> New Ticket</Button>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New Ticket</DialogTitle>
        <DialogDescription>
          Create a new {{ props.selectedType.singular }} ticket.
        </DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit">
        <FormField name="name" v-slot="{ componentField }" v-model="name">
          <FormItem>
            <FormLabel for="name" class="text-right">Name</FormLabel>
            <Input id="name" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="description" v-slot="{ componentField }" v-model="description">
          <FormItem>
            <FormLabel for="description" class="text-right">Description</FormLabel>
            <Input id="description" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <JSONSchemaFormFields v-model="state" :schema="selectedType.schema" />

        <DialogFooter class="mt-4 sm:justify-start">
          <Button type="submit"> Save </Button>
          <DialogClose as-child>
            <Button type="button" variant="secondary">Cancel</Button>
          </DialogClose>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
