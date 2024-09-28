<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { onMounted, ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Resource, Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const isOpen = defineModel<boolean>()

const addResourceMutation = useMutation({
  mutationFn: (values: any): Promise<Resource> => {
    return pb.collection('resources').create({
      ticket: props.ticket.id,
      value: values.value
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    isOpen.value = false
  },
  onError: handleError
})

/*
const isURL = (value: string) => {
  try {
    new URL(value)
    return true
  } catch {
    return false
  }
}

const addResourceMutation = useMutation({
  mutationFn: (values: any): Promise<Resource> => {
    return pb
      .send('/api/analysis/enrich', { value: values.value })
      .then((data) => {
        const resources = data.resources
        if (resources.length > 0) {
          return pb.collection('resources').create({
            ticket: props.ticket.id,
            service: resources[0].service,
            type: resources[0].type,
            resource: resources[0].id,
            name: resources[0].name,
            icon: resources[0].icon,
            description: resources[0].description,
            url: resources[0].url,
            attributes: resources[0].attributes
          })
        }

        return pb.collection('resources').create({
          ticket: props.ticket.id,
          service: 'catalyst',
          type: 'other',
          resource: values.value,
          name: values.value,
          icon: 'Link',
          description: '',
          url: isURL(values.value) ? values.value : '',
          attributes: []
        })
      })
      .catch((err) => {
        return Promise.reject({
          name: "Error",
          message: "An error occurred" + JSON.stringify(err)
        })
      })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    isOpen.value = false
  },
  onError: handleError
})
 */

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const { handleSubmit, validate } = useForm({
  validationSchema: {
    value: 'required'
  }
})

const onSubmit = handleSubmit((values) => addResourceMutation.mutate(values))

const submitDisabled = ref(true)
onMounted(() => change())

const change = () => validate({ mode: 'silent' }).then((res) => (submitDisabled.value = !res.valid))
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New Resource</DialogTitle>
        <DialogDescription> Add a new resource to this ticket</DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit" @change="change">
        <FormField name="value" v-slot="{ componentField }" class="mt-2">
          <FormItem>
            <FormLabel for="value" class="text-right">Value</FormLabel>
            <Input id="value" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter class="mt-2 sm:justify-start">
          <Button
            :title="submitDisabled ? 'Please fill out all required fields' : undefined"
            :disabled="submitDisabled"
            type="submit"
          >
            Save
          </Button>
          <DialogClose as-child>
            <Button type="button" variant="secondary"> Cancel</Button>
          </DialogClose>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
