<script setup lang="ts">
import { Button, buttonVariants } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import { CircleUser } from 'lucide-vue-next'

import type { AuthModel } from 'pocketbase'
import { onMounted, ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import { cn } from '@/lib/utils'

defineProps<{
  isCollapsed: boolean
}>()

const variant = 'secondary'

interface User {
  name: string
}

const user = ref<AuthModel | User>(pb.authStore.model)

const logout = () => {
  pb.authStore.clear()
  window.location.href = '/ui/login'
}

onMounted(() => {
  pb.collection('users')
    .authRefresh()
    .catch(() => {
      pb.authStore.clear()
      window.location.href = '/ui/login'
    })
})
</script>

<template>
  <div class="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2">
    <nav
      v-if="user"
      class="grid gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2"
    >
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <div>
            <Tooltip v-if="isCollapsed" :delay-duration="0">
              <TooltipTrigger as-child>
                <Button
                  :class="
                    cn(buttonVariants({ variant: variant, size: 'icon' }), 'mx-1 h-9 w-9 px-0')
                  "
                >
                  <CircleUser class="size-4" />
                  <span class="sr-only">{{ user.name }}</span>
                </Button>
              </TooltipTrigger>
              <TooltipContent side="right" class="flex items-center gap-4">
                {{ user.name }}
              </TooltipContent>
            </Tooltip>
            <Button
              v-else
              :class="cn(buttonVariants({ variant: variant, size: 'sm' }), 'w-full justify-start')"
            >
              <CircleUser class="mr-2 size-4" />
              {{ user.name }}
            </Button>
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuLabel>Account</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            @click="logout"
            class="cursor-pointer text-muted-foreground transition-colors hover:text-foreground"
          >
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </nav>
  </div>
</template>
