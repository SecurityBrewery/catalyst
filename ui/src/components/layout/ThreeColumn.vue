<script lang="ts" setup>
import SideBar from '@/components/layout/SideBar.vue'
import { TooltipProvider } from '@/components/ui/tooltip'

import { cn } from '@/lib/utils'

defineProps<{
  showDetails?: boolean
}>()
</script>

<template>
  <TooltipProvider :delay-duration="0">
    <div class="flex h-full flex-row items-stretch bg-muted/40">
      <SideBar />
      <div
        :class="
          cn(
            'w-full flex-initial border-r sm:w-72',
            !showDetails && 'flex',
            showDetails && 'hidden sm:flex'
          )
        "
      >
        <div class="flex h-full w-full flex-col">
          <slot name="list" />
        </div>
      </div>
      <div
        :class="
          cn('flex-1 overflow-hidden', !showDetails && 'hidden sm:flex', showDetails && 'flex')
        "
      >
        <div class="flex h-full w-full flex-1 flex-col">
          <slot name="single" />
        </div>
      </div>
    </div>
  </TooltipProvider>
</template>
