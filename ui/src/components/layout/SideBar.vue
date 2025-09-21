<script setup lang="ts">
import Icon from '@/components/Icon.vue'
import CatalystLogo from '@/components/common/CatalystLogo.vue'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail
} from '@/components/ui/sidebar'

import { ChevronsUpDown, LogOut, Settings, Tag, User, Users, Zap } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api.ts'
import type { Sidebar as SidebarModel } from '@/client/models'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()

const links = [
  {
    name: 'Reactions',
    url: '/reactions',
    icon: Zap,
    permission: 'reaction:write'
  },
  {
    name: 'Users',
    url: '/users',
    icon: User,
    permission: 'group:write'
  },
  {
    name: 'Groups',
    url: '/groups',
    icon: Users,
    permission: 'group:write'
  },
  {
    name: 'Types',
    url: '/types',
    icon: Tag,
    permission: 'type:write'
  },
  {
    name: 'Settings',
    url: '/settings',
    icon: Settings,
    permission: 'settings:write'
  }
]

const userLinks = computed(() => {
  return links.filter((link) => authStore.hasPermission(link.permission))
})

const api = useAPI()
const router = useRouter()

const { data: sidebar } = useQuery({
  queryKey: ['sidebar'],
  queryFn: (): Promise<Array<SidebarModel>> => api.getSidebar()
})

const logout = () => {
  authStore.setToken('')
  router.push({ name: 'login' })
}

const initials = (user: { name?: string } | undefined) => {
  if (!user || !user.name) return ''
  const names = user.name.split(' ')
  return names.length > 1 ? `${names[0][0]}${names[1][0]}` : names[0][0]
}
</script>

<template>
  <SidebarProvider class="h-full w-full" style="--sidebar-width: 12rem">
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem class="flex items-center gap-2 p-2">
            <CatalystLogo :size="5" />
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">Catalyst</span>
              <span class="truncate text-xs">Incident Management</span>
            </div>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Overview</SidebarGroupLabel>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton :tooltip="'Dashboard'" as-child>
                <RouterLink to="/dashboard">
                  <Icon name="LayoutDashboard" class="size-4" />
                  <span>Dashboard</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Tickets</SidebarGroupLabel>
          <SidebarMenu>
            <SidebarMenuItem v-for="(typ, index) of sidebar" :key="index">
              <SidebarMenuButton :tooltip="typ.plural" as-child>
                <RouterLink :to="`/tickets/${typ.id}`">
                  <Icon :name="typ.icon" class="size-4" />
                  <span>{{ typ.plural }}</span>
                  <span class="ml-auto">{{ typ.count }}</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
        <div class="flex-1" />
        <SidebarGroup v-if="userLinks.length > 0" class="group-data-[collapsible=icon]:hidden">
          <SidebarGroupLabel>Administration</SidebarGroupLabel>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in userLinks" :key="item.name">
              <SidebarMenuButton as-child>
                <RouterLink :to="item.url">
                  <component :is="item.icon" />
                  <span>{{ item.name }}</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <SidebarMenuButton
                  size="lg"
                  class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <Avatar class="h-8 w-8 rounded-lg">
                    <AvatarImage
                      :src="authStore.user?.avatar ? authStore.user.avatar : ''"
                      :alt="authStore.user?.name"
                    />
                    <AvatarFallback class="rounded-lg"
                      >{{ initials(authStore.user) }}
                    </AvatarFallback>
                  </Avatar>
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate font-semibold">{{ authStore.user?.name }}</span>
                    <span class="truncate text-xs">{{ authStore.user?.email }}</span>
                  </div>
                  <ChevronsUpDown class="ml-auto size-4" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                class="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                side="bottom"
                align="end"
                :side-offset="4"
              >
                <DropdownMenuLabel class="p-0 font-normal">
                  <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <Avatar class="h-8 w-8 rounded-lg">
                      <AvatarImage
                        :src="authStore.user?.avatar ? authStore.user.avatar : ''"
                        :alt="authStore.user?.name"
                      />
                      <AvatarFallback class="rounded-lg"
                        >{{ initials(authStore.user) }}
                      </AvatarFallback>
                    </Avatar>
                    <div class="grid flex-1 text-left text-sm leading-tight">
                      <span class="truncate font-semibold">{{ authStore.user?.name }}</span>
                      <span class="truncate text-xs">{{ authStore.user?.email }}</span>
                    </div>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem @click="logout">
                  <LogOut />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
    <SidebarInset class="flex h-full w-full">
      <slot />
    </SidebarInset>
  </SidebarProvider>
</template>
