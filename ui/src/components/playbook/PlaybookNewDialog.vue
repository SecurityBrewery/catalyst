<script setup lang="ts">
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
import { toast } from '@/components/ui/toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Playbook, Ticket } from '@/lib/types'

const queryClient = useQueryClient()
const router = useRouter()

const isOpen = ref(false)

const addPlaybookMutation = useMutation({
  mutationFn: (values: any): Promise<Playbook> =>
    pb.collection('playbooks').create({
      name: values.name,
      schema: {
        steps: []
      }
    }),
  onSuccess: (data: Ticket) => {
    router.push(`/playbooks/${data.id}`)
    queryClient.invalidateQueries({ queryKey: ['playbooks'] })
    isOpen.value = false
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const { handleSubmit, validate } = useForm({
  validationSchema: {
    name: 'required'
  }
})

const onSubmit = handleSubmit((values) => addPlaybookMutation.mutate(values))

const submitDisabled = ref(true)
onMounted(() => change())

const change = () => validate({ mode: 'silent' }).then((res) => (submitDisabled.value = !res.valid))
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="ghost"> New Playbook</Button>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New Playbook</DialogTitle>
        <DialogDescription> Create a new playbook</DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit" @change="change">
        <FormField name="name" v-slot="{ componentField }">
          <FormItem>
            <FormLabel for="name" class="text-right"> Name</FormLabel>
            <Input id="name" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter class="mt-2">
          <DialogClose as-child>
            <Button type="button" variant="secondary"> Cancel</Button>
          </DialogClose>
          <Button
            :title="submitDisabled ? 'Please fill out all required fields' : undefined"
            :disabled="submitDisabled"
            type="submit"
            >Save
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
