<script setup lang="ts">
import { Separator } from '@/components/ui/separator'
import { SidebarTrigger } from '@/components/ui/sidebar'

import { cn } from '@/lib/utils'

defineProps<{
  showSidebarTrigger?: boolean
  title?: string
  nowrap?: boolean
  hideSeparator?: boolean
}>()
</script>

<template>
  <div
    :class="
      cn('flex min-h-14 flex-wrap items-center gap-2 bg-background p-2', nowrap && 'flex-nowrap')
    "
  >
    <SidebarTrigger v-if="showSidebarTrigger" />
    <h1 v-if="title" class="text-xl font-bold">
      {{ title }}
    </h1>
    <slot />
  </div>
  <Separator v-if="!hideSeparator" />
</template>
