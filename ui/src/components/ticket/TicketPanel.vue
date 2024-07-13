<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

import { Plus } from 'lucide-vue-next'

const emit = defineEmits(['add'])

export interface Props {
  title: string
  hideAdd?: boolean
}

withDefaults(defineProps<Props>(), {
  hideAdd: false
})
</script>

<template>
  <div>
    <div class="flex h-10 flex-row items-center justify-between text-muted-foreground">
      <span class="text-sm font-semibold">
        {{ title }}
      </span>

      <Button v-if="!hideAdd" variant="ghost" size="icon" class="h-8 w-8" @click="emit('add')">
        <Plus class="size-4" />
        <span class="sr-only">Add item</span>
      </Button>
    </div>
    <Card v-if="$slots.default" class="p-0">
      <slot />
    </Card>
  </div>
</template>
