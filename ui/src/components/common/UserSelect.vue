<script setup lang="ts">
import UserSelectList from '@/components/common/UserSelectList.vue'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

import { ref } from 'vue'

import type { User } from '@/lib/types'

const user = defineModel<User>()

const open = ref(false)
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <slot />
    </PopoverTrigger>
    <PopoverContent class="w-[150px] p-0">
      <UserSelectList v-model="user" :key="user ? user.id : 'unassigned'" :user="user" />
    </PopoverContent>
  </Popover>
</template>
