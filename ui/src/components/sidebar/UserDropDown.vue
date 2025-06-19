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

import { useRouter } from 'vue-router'

import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

defineProps<{
  isCollapsed: boolean
}>()

const variant = 'secondary'

const authStore = useAuthStore()
const router = useRouter()

const logout = () => {
  authStore.setToken('')
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2">
    <nav
      v-if="authStore.user"
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
                  <span class="sr-only">{{ authStore.user.name }}</span>
                </Button>
              </TooltipTrigger>
              <TooltipContent side="right" class="flex items-center gap-4">
                {{ authStore.user.name }}
              </TooltipContent>
            </Tooltip>
            <Button
              v-else
              :class="cn(buttonVariants({ variant: variant, size: 'sm' }), 'w-full justify-start')"
            >
              <CircleUser class="mr-2 size-4" />
              {{ authStore.user.name }}
            </Button>
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuLabel>Account</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            @click="logout"
            class="text-muted-foreground hover:text-foreground cursor-pointer transition-colors"
          >
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </nav>
  </div>
</template>
