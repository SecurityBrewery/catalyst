<script setup lang="ts">
import { Separator } from '@/components/ui/separator'
import { TabsContent } from '@/components/ui/tabs'

defineProps<{
  value: string
}>()
</script>

<template>
  <TabsContent :value="value" class="flex-1 overflow-hidden">
    <div class="flex h-full flex-col overflow-hidden">
      <Separator class="mt-2" />
      <slot />
      <div class="h-4" />
    </div>
  </TabsContent>
</template>
