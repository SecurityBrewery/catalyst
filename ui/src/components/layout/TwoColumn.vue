<script lang="ts" setup>
import SideBar from '@/components/layout/SideBar.vue'
import { TooltipProvider } from '@/components/ui/tooltip'
</script>

<template>
  <TooltipProvider :delay-duration="0">
    <div class="bg-muted/40 flex h-full flex-row items-stretch">
      <SideBar />
      <div class="flex h-full w-full flex-col">
        <slot />
      </div>
    </div>
  </TooltipProvider>
</template>
