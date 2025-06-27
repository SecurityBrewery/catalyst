<script lang="ts" setup>
import SideBar from '@/components/layout/SideBar.vue'
import { TooltipProvider } from '@/components/ui/tooltip'
</script>

<template>
  <TooltipProvider :delay-duration="0">
    <SideBar>
      <slot />
    </SideBar>
  </TooltipProvider>
</template>
