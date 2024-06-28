<script lang="ts" setup>
import Icon from '@/components/Icon.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { buttonVariants } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import { LoaderCircle } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import { useRoute } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Type } from '@/lib/types'
import { cn } from '@/lib/utils'

const route = useRoute()

defineProps<{
  isCollapsed: boolean
}>()

const {
  isPending,
  isError,
  data: sidebar,
  error
} = useQuery({
  queryKey: ['sidebar'],
  queryFn: (): Promise<Array<any>> => pb.collection('sidebar').getFullList()
})

const variant = (t: Type): 'default' | 'ghost' => (route.params.type === t.id ? 'default' : 'ghost')
</script>

<template>
  <div
    :data-collapsed="isCollapsed"
    class="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2"
  >
    <div v-if="isPending" class="flex h-screen w-screen items-center justify-center">
      <LoaderCircle class="h-16 w-16 animate-spin text-primary" />
    </div>
    <Alert v-else-if="isError" variant="destructive" class="mb-4 h-screen w-screen">
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <nav
      v-else-if="sidebar"
      class="grid gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2"
    >
      <template v-for="(typ, index) of sidebar">
        <Tooltip v-if="isCollapsed" :key="`1-${index}`" :delay-duration="0">
          <TooltipTrigger as-child>
            <RouterLink
              :to="`/tickets/${typ.id}`"
              :class="
                cn(
                  buttonVariants({ variant: variant(typ), size: 'icon' }),
                  'h-9 w-9',
                  variant(typ) === 'default' &&
                    'dark:bg-muted dark:text-muted-foreground dark:hover:bg-muted dark:hover:text-white'
                )
              "
            >
              <Icon :name="typ.icon" class="size-4" />
              <span class="sr-only">{{ typ.plural }}</span>
            </RouterLink>
          </TooltipTrigger>
          <TooltipContent side="right" class="flex items-center gap-4">
            {{ typ.plural }}
            <span v-if="typ.count" class="ml-auto text-muted-foreground">
              {{ typ.count }}
            </span>
          </TooltipContent>
        </Tooltip>

        <RouterLink
          v-else
          :key="`2-${index}`"
          :to="`/tickets/${typ.id}`"
          :class="
            cn(
              buttonVariants({ variant: variant(typ), size: 'sm' }),
              variant(typ) === 'default' &&
                'dark:bg-muted dark:text-white dark:hover:bg-muted dark:hover:text-white',
              'justify-start'
            )
          "
        >
          <Icon :name="typ.icon" class="mr-2 size-4" />
          {{ typ.plural }}
          <span
            v-if="typ.count"
            :class="cn('ml-auto', variant(typ) === 'default' && 'text-background dark:text-white')"
          >
            {{ typ.count }}
          </span>
        </RouterLink>
      </template>
    </nav>
  </div>
</template>
