<script setup lang="ts">
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '@/components/ui/command'

import { Check } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import debounce from 'lodash.debounce'
import { ref, watch } from 'vue'

import { api } from '@/api'
import type { User } from '@/client'
import { cn } from '@/lib/utils'

defineProps<{
  userID: string | undefined
}>()

const searchTerm = ref('')

const {
  isPending: usersIsPending,
  isError: usersIsError,
  data: users,
  error: usersError,
  refetch
} = useQuery({
  queryKey: ['users', 'search', searchTerm.value],
  queryFn: () =>
    api
      .listUsers()
      .then((users) =>
        users.filter(
          (user) =>
            user.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
            user.username.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
            user.email.toLowerCase().includes(searchTerm.value.toLowerCase())
        )
      )
})

const searchUserDebounced = debounce(() => refetch(), 300)

watch(
  () => searchTerm.value,
  () => searchUserDebounced()
)

const emit = defineEmits<{
  (e: 'select', value: User): void
}>()
</script>

<template>
  <Command v-model:search-term="searchTerm">
    <CommandInput placeholder="Search user..." />
    <CommandEmpty>
      <span v-if="usersIsPending"> Loading... </span>
      <span v-else-if="usersIsError"> Error: {{ usersError }} </span>
      <span>No user found.</span>
    </CommandEmpty>
    <CommandList>
      <CommandGroup>
        <CommandItem
          v-for="u in users"
          :key="u.id"
          :value="u"
          @select="emit('select', u)"
          class="cursor-pointer"
        >
          <Check
            :class="cn('mr-2 h-4 w-4', userID && userID === u.id ? 'opacity-100' : 'opacity-0')"
          />
          {{ u.name }}
        </CommandItem>
      </CommandGroup>
    </CommandList>
  </Command>
</template>
