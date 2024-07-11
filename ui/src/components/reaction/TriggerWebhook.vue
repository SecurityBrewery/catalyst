<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

import { computed } from 'vue'

interface TriggerWebhookData {
  token: string
  path: string
}

const props = defineProps<{
  modelValue: TriggerWebhookData
}>()

const emit = defineEmits(['update:modelValue'])

const updateToken = (value: string) =>
  emit('update:modelValue', { ...props.modelValue, token: value })
const updatePath = (value: string) =>
  emit('update:modelValue', { ...props.modelValue, path: value })

const curlExample = computed(() => {
  let cmd = `curl`

  if (props.modelValue && props.modelValue.token) {
    cmd += ` -H "Auth: Bearer ${props.modelValue.token}"`
  }

  if (props.modelValue && props.modelValue.path) {
    cmd += ` https://${location.hostname}/reaction/${props.modelValue.path}`
  }

  return cmd
})
</script>

<template>
  <p class="py-4 text-sm text-muted-foreground">
    Receive a POST request at the specified path with an optional authorization token.
  </p>

  <Label for="token">Token</Label>
  <Input
    id="token"
    :modelValue="modelValue.token"
    @update:modelValue="updateToken"
    placeholder="Enter a token (e.g. 'Bearer xyz...')"
  />

  <Label for="path">Path</Label>
  <Input
    id="path"
    :modelValue="modelValue.path"
    @update:modelValue="updatePath"
    placeholder="Enter a path (e.g. '/webhook')"
  />

  <Label for="url">Usage</Label>
  <Input id="url" readonly :modelValue="curlExample" class="bg-accent" />
</template>
