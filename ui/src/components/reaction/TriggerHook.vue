<script setup lang="ts">
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

interface TriggerWebhookData {
  collection: string
  create: boolean
  update: boolean
  delete: boolean
}

const props = defineProps<{
  modelValue: TriggerWebhookData
}>()

const emit = defineEmits(['update:modelValue'])

const updateCollection = (value: string) =>
  emit('update:modelValue', { ...props.modelValue, collection: value })
const updateCreate = (value: boolean) =>
  emit('update:modelValue', { ...props.modelValue, create: value })
const updateUpdate = (value: boolean) =>
  emit('update:modelValue', { ...props.modelValue, update: value })
const updateDelete = (value: boolean) =>
  emit('update:modelValue', { ...props.modelValue, delete: value })
</script>

<template>
  <p class="py-4 text-sm text-muted-foreground">
    Trigger an reaction when an event occurs in a collection.
  </p>

  <Label for="collection">Collection</Label>
  <Select id="collection" :modelValue="modelValue.collection" @update:modelValue="updateCollection">
    <SelectTrigger>
      <SelectValue placeholder="Select a collection" />
    </SelectTrigger>
    <SelectContent>
      <SelectGroup>
        <SelectItem value="tickets">Tickets</SelectItem>
        <SelectItem value="tasks">Tasks</SelectItem>
        <SelectItem value="comments">Comments</SelectItem>
        <SelectItem value="timeline">Timeline</SelectItem>
        <SelectItem value="links">Links</SelectItem>
        <SelectItem value="files">Files</SelectItem>
      </SelectGroup>
    </SelectContent>
  </Select>

  <div class="mt-2 flex items-center space-x-2">
    <Switch id="create" :checked="modelValue.create" @update:checked="updateCreate" />
    <Label for="create">Create events</Label>
  </div>

  <div class="mt-2 flex items-center space-x-2">
    <Switch id="update" :checked="modelValue.update" @update:checked="updateUpdate" />
    <Label for="update">Update events</Label>
  </div>

  <div class="mt-2 flex items-center space-x-2">
    <Switch id="delete" :checked="modelValue.delete" @update:checked="updateDelete" />
    <Label for="delete">Delete events</Label>
  </div>
</template>
