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
    <SideBar>
      <div class="flex h-full w-full flex-col overflow-auto md:flex-row">
        <div :class="cn('flex h-full w-full flex-col border-r md:w-72', showDetails ? 'hidden md:flex' : '')">
          <slot name="list" />
        </div>
        <div :class="cn('flex h-full w-full flex-1 flex-col', showDetails ? '' : 'hidden md:flex')">
          <slot name="single" />
        </div>
      </div>
    </SideBar>
  </TooltipProvider>
</template>
