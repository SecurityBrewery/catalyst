<script setup lang="ts">
import { cn } from '@/lib/utils'
import { useCatalystStore } from '@/store/catalyst'

const catalystStore = useCatalystStore()

defineProps<{
  size: number
}>()
</script>

<template>
  <img
    src="@/assets/flask.svg"
    alt="Catalyst Logo"
    :class="
      cn('dark:hidden', 'size-' + size, {
        'flex-1': catalystStore.sidebarCollapsed,
        'mx-3': !catalystStore.sidebarCollapsed
      })
    "
  />
  <img
    src="@/assets/flask_white.svg"
    alt="Catalyst Logo"
    :class="
      cn('hidden dark:flex', 'size-' + size, {
        'flex-1': catalystStore.sidebarCollapsed,
        'mx-3': !catalystStore.sidebarCollapsed
      })
    "
  />
</template>
