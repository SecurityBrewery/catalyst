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

import { pb } from '@/lib/pocketbase'
import type { User } from '@/lib/types'
import { cn } from '@/lib/utils'

const user = defineModel<User>()

const open = ref(false)
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
    pb.collection('users').getFullList({
      sort: 'name',
      perPage: 5,
      filter: pb.filter(`name ~ {:search}`, { search: searchTerm.value })
    })
})

const searchUserDebounced = debounce(() => refetch(), 300)

watch(
  () => searchTerm.value,
  () => searchUserDebounced()
)
</script>

<template>
  <Command v-model="user" v-model:search-term="searchTerm">
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
          @select="open = false"
          class="cursor-pointer"
        >
          <Check
            :class="cn('mr-2 h-4 w-4', user && user.id === u.id ? 'opacity-100' : 'opacity-0')"
          />
          {{ u.name }}
        </CommandItem>
      </CommandGroup>
    </CommandList>
  </Command>
</template>
