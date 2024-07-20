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
import type { Link, Ticket } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const isOpen = defineModel<boolean>()

const addLinkMutation = useMutation({
  mutationFn: (values: any): Promise<Link> =>
    pb.collection('links').create({
      ticket: props.ticket.id,
      name: values.name,
      url: values.url
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
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

const { handleSubmit, validate } = useForm({
  validationSchema: {
    name: 'required',
    url: 'required'
  }
})

const onSubmit = handleSubmit((values) => addLinkMutation.mutate(values))

const submitDisabled = ref(true)
onMounted(() => change())

const change = () => validate({ mode: 'silent' }).then((res) => (submitDisabled.value = !res.valid))
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New Link</DialogTitle>
        <DialogDescription> Add a new link to this ticket</DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit" @change="change">
        <FormField name="name" v-slot="{ componentField }">
          <FormItem>
            <FormLabel for="name" class="text-right"> Name</FormLabel>
            <Input id="name" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="url" v-slot="{ componentField }" class="mt-2">
          <FormItem>
            <FormLabel for="url" class="text-right"> URL</FormLabel>
            <Input id="url" v-bind="componentField" />
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
