<script setup lang="ts">
import CatalystLogo from '@/components/common/CatalystLogo.vue'
import IncidentNav from '@/components/sidebar/IncidentNav.vue'
import NavList from '@/components/sidebar/NavList.vue'
import UserDropDown from '@/components/sidebar/UserDropDown.vue'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'

import { Menu } from 'lucide-vue-next'

import { cn } from '@/lib/utils'
import { useCatalystStore } from '@/store/catalyst'
import { useAuthStore } from '@/store/auth'

const catalystStore = useCatalystStore()
const authStore = useAuthStore()
</script>

<template>
  <div
    :class="
      cn(
        'flex min-w-48 shrink-0 flex-col border-r bg-popover', // transition-all duration-300 ease-in-out',
        catalystStore.sidebarCollapsed && 'min-w-[50px]'
      )
    "
  >
    <div class="flex h-[57px] items-center border-b bg-background">
      <CatalystLogo :size="8"/>
      <h1 class="text-xl font-bold" v-if="!catalystStore.sidebarCollapsed">Catalyst</h1>
    </div>
    <div v-if="!catalystStore.sidebarCollapsed && !authStore.permissions.includes('ticket:read')" class="w-full px-2 mt-4">
      <Alert class="w-full">
        <AlertTitle>Info</AlertTitle>
        <AlertDescription>No permission to read tickets</AlertDescription>
      </Alert>
    </div>
    <NavList
      :is-collapsed="catalystStore.sidebarCollapsed"
      :links="[
        {
          title: 'Dashboard',
          icon: 'PanelsTopLeft',
          variant: 'ghost',
          to: '/dashboard',
          permission: 'ticket:read'
        }
      ]"
    />
    <Separator v-if="authStore.permissions.includes('ticket:read')" />
    <IncidentNav :is-collapsed="catalystStore.sidebarCollapsed" />

    <div class="flex-1" />

    <Separator v-if="authStore.permissions.includes('reaction:write') || authStore.permissions.includes('role:write') || authStore.permissions.includes('group:write') || authStore.permissions.includes('type:write')" />
    <NavList
      :is-collapsed="catalystStore.sidebarCollapsed"
      :links="[
        {
          title: 'Reactions',
          icon: 'Zap',
          variant: 'ghost',
          to: '/reactions',
          permission: 'reaction:write'
        },
        {
          title: 'Users',
          icon: 'User',
          variant: 'ghost',
          to: '/users',
          permission: 'role:write'
        },
        {
          title: 'Groups',
          icon: 'Users',
          variant: 'ghost',
          to: '/groups',
          permission: 'group:write'
        },
        {
          title: 'Types',
          icon: 'Tag',
          variant: 'ghost',
          to: '/types',
          permission: 'type:write'
        }
      ]"
    />
    <Separator />
    <UserDropDown :is-collapsed="catalystStore.sidebarCollapsed" />
    <Separator />
    <div :class="cn('flex h-14 items-center px-3', !catalystStore.sidebarCollapsed && 'px-2')">
      <Button
        variant="ghost"
        @click="catalystStore.toggleSidebar()"
        size="default"
        :class="
          cn(
            'p-0',
            catalystStore.sidebarCollapsed && 'w-9',
            !catalystStore.sidebarCollapsed && 'w-full justify-start px-3'
          )
        "
      >
        <Menu class="size-4" />
        <span v-if="!catalystStore.sidebarCollapsed" class="ml-2">Toggle Sidebar</span>
      </Button>
    </div>
  </div>
</template>
