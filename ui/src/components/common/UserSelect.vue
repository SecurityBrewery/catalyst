<script setup lang="ts">
import UserSelectList from '@/components/common/UserSelectList.vue'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

import { ref } from 'vue'

const userID = defineModel<string | undefined>()

const open = ref(false)
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <slot />
    </PopoverTrigger>
    <PopoverContent class="w-[150px] p-0">
      <UserSelectList v-model="userID" :key="userID ? userID : 'unassigned'" :userID="userID" />
    </PopoverContent>
  </Popover>
</template>
