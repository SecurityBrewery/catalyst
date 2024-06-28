<script setup lang="ts">
import { ScrollArea } from '@/components/ui/scroll-area'
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
      <ScrollArea class="flex-1">
        <slot />
        <div class="h-4" />
      </ScrollArea>
    </div>
  </TabsContent>
</template>
