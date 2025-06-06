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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'

import { useQuery } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { User } from '@/client'

const api = useAPI()

const isOpen = defineModel<boolean>()

const props = defineProps<{
  exclude: Array<string>
}>()

const emit = defineEmits(['select'])

const { data: users } = useQuery({
  queryKey: ['users'],
  queryFn: (): Promise<Array<User>> => api.listUsers()
})

const filteredUsers = computed(() => {
  return users.value?.filter((user) => !props.exclude.includes(user.id)) ?? []
})

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const { handleSubmit, validate, values } = useForm({
  validationSchema: {
    user: 'required'
  }
})

const onSubmit = handleSubmit((values) => emit('select', values))

const submitDisabled = ref(true)

const change = () => validate({ mode: 'silent' }).then((res) => (submitDisabled.value = !res.valid))

watch(
  () => values,
  () => change(),
  { deep: true, immediate: true }
)
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New User</DialogTitle>
        <DialogDescription> Add a new user to this group</DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit" @change="change">
        <FormField name="user" v-slot="{ componentField }">
          <FormItem>
            <FormLabel for="user" class="text-right"> User</FormLabel>
            <Select id="user" v-bind="componentField">
              <SelectTrigger>
                <SelectValue placeholder="Select a user" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="user in filteredUsers" :key="user.id" :value="user.id">
                  {{ user.name }}
                </SelectItem>
              </SelectContent>
            </Select>
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
