<script setup lang="ts">
import CatalystLogo from '@/components/common/CatalystLogo.vue'
import IncidentNav from '@/components/sidebar/IncidentNav.vue'
import NavList from '@/components/sidebar/NavList.vue'
import UserDropDown from '@/components/sidebar/UserDropDown.vue'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

import { Menu } from 'lucide-vue-next'

import { useCatalystStore } from '@/store/catalyst'

const catalystStore = useCatalystStore()
</script>

<template>
  <div class="flex h-[57px] items-center border-b bg-background">
    <CatalystLogo
      class="size-8"
      :class="{ 'flex-1': catalystStore.sidebarCollapsed, 'mx-3': !catalystStore.sidebarCollapsed }"
    />
    <h1 class="text-xl font-bold" v-if="!catalystStore.sidebarCollapsed">Catalyst</h1>
  </div>
  <NavList
    :is-collapsed="catalystStore.sidebarCollapsed"
    :links="[
      {
        title: 'Dashboard',
        icon: 'PanelsTopLeft',
        variant: 'ghost',
        to: '/dashboard'
      }
    ]"
  />
  <Separator />
  <IncidentNav :is-collapsed="catalystStore.sidebarCollapsed" />

  <div class="flex-1" />

  <Separator />
  <NavList
    :is-collapsed="catalystStore.sidebarCollapsed"
    :links="[
      {
        title: 'Reactions',
        icon: 'Zap',
        variant: 'ghost',
        to: '/reactions'
      }
    ]"
  />
  <Separator />
  <UserDropDown :is-collapsed="catalystStore.sidebarCollapsed" />
  <Separator />
  <Button
    variant="ghost"
    @click="catalystStore.toggleSidebar()"
    size="sm"
    class="m-2 justify-start px-3.5"
  >
    <Menu class="size-4" />
    <span v-if="!catalystStore.sidebarCollapsed" class="ml-2">Toggle Sidebar</span>
  </Button>
</template>
