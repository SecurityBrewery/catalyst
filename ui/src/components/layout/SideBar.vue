<script setup lang="ts">
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
    <img
      src="@/assets/flask.svg"
      alt="Catalyst"
      class="h-8 w-8 dark:hidden"
      :class="{ 'flex-1': catalystStore.sidebarCollapsed, 'mx-3': !catalystStore.sidebarCollapsed }"
    />
    <img
      src="@/assets/flask_white.svg"
      alt="Catalyst"
      class="hidden h-8 w-8 dark:flex"
      :class="{ 'flex-1': catalystStore.sidebarCollapsed, 'mx-3': !catalystStore.sidebarCollapsed }"
    />
    <h1 class="text-xl font-bold" v-if="!catalystStore.sidebarCollapsed">Catalyst</h1>
  </div>
  <NavList
    class="mt-auto"
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

  <Separator />
  <NavList
    :is-collapsed="catalystStore.sidebarCollapsed"
    :links="[
      {
        title: 'Playbooks',
        icon: 'Book',
        variant: 'ghost',
        to: '/playbooks'
      }
    ]"
  />
  <Separator />

  <div class="flex-1" />

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
