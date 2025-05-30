<script setup lang="ts">
import UserSelectList from '@/components/common/UserSelectList.vue'
import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

import { User2 } from 'lucide-vue-next'

import type { User } from '@/client'

defineProps<{
  userID?: string
  userName?: string
}>()

const select = (user: User) => {
  emit('select', user)
}

const emit = defineEmits<{
  (e: 'select', value: User): void
}>()
</script>

<template>
  <Popover>
    <PopoverTrigger as-child>
      <Button variant="outline" role="combobox">
        <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
        {{ userName ? userName : 'Unassigned' }}
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-[150px] p-0">
      <UserSelectList :userID="userID" @select="select" />
    </PopoverContent>
  </Popover>
</template>
