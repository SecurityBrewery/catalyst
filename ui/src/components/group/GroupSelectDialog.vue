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
import type { Group } from '@/client'

const api = useAPI()

const isOpen = defineModel<boolean>()

const props = defineProps<{
  exclude: Array<string>
}>()

const emit = defineEmits(['select'])

const { data: groups } = useQuery({
  queryKey: ['groups'],
  queryFn: (): Promise<Array<Group>> => api.listGroups()
})

const filteredGroups = computed(() => {
  return groups.value?.filter((group) => !props.exclude.includes(group.id)) ?? []
})

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const { handleSubmit, validate, values } = useForm({
  validationSchema: {
    group: 'required'
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
        <DialogTitle>New Group</DialogTitle>
        <DialogDescription> Add a new group to this user</DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit" @change="change">
        <FormField name="group" v-slot="{ componentField }">
          <FormItem>
            <FormLabel for="group" class="text-left"> Group</FormLabel>
            <Select id="group" v-bind="componentField">
              <SelectTrigger>
                <SelectValue placeholder="Select a group" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="group in filteredGroups" :key="group.id" :value="group.id"
                  >{{ group.name }}
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
