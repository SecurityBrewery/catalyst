<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

import { LoaderCircle } from 'lucide-vue-next'

defineProps<{
  isPending: boolean
  isError: boolean
  error: Error | null
}>()
</script>

<template>
  <div v-if="isPending" class="flex h-full w-full">
    <LoaderCircle class="m-auto h-16 w-16 animate-spin text-primary" />
  </div>
  <Alert v-else-if="isError" variant="destructive" class="mb-4">
    <AlertTitle>Error</AlertTitle>
    <AlertDescription>{{ error }}</AlertDescription>
  </Alert>
  <slot v-else />
</template>
