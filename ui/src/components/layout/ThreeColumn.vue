<script lang="ts" setup>
import SideBar from '@/components/layout/SideBar.vue'
import { TooltipProvider } from '@/components/ui/tooltip'

import { cn } from '@/lib/utils'
import { useCatalystStore } from '@/store/catalyst'

const catalystStore = useCatalystStore()
</script>

<template>
  <TooltipProvider :delay-duration="0">
    <div class="flex h-full flex-row items-stretch bg-muted/40">
      <div
        :class="
          cn(
            'flex min-w-48 flex-col border-r bg-popover', // transition-all duration-300 ease-in-out',
            catalystStore.sidebarCollapsed && 'min-w-[50px]'
          )
        "
      >
        <SideBar />
      </div>
      <div class="w-72 flex-initial border-r">
        <slot name="list" />
      </div>
      <div class="flex-1">
        <slot name="single" />
      </div>
    </div>
  </TooltipProvider>
</template>
