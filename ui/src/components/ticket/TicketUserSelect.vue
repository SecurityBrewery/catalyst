<script setup lang="ts">
import UserSelect from '@/components/common/UserSelect.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'

import { LoaderCircle, User2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Ticket, User } from '@/lib/types'
import { handleError } from '@/lib/utils'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  uID?: string
}>()

const {
  isPending,
  isError,
  data: user,
  error
} = useQuery({
  queryKey: ['tickets', props.ticket.id, 'owner', props.uID],
  queryFn: (): Promise<User | null> => {
    if (!props.uID) {
      return Promise.resolve(null)
    }

    return pb.collection('users').getOne(props.uID)
  }
})

const setTicketOwnerMutation = useMutation({
  mutationFn: (user: User): Promise<Ticket> =>
    pb.collection('tickets').update(props.ticket.id, {
      owner: user.id
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets'] }),
  onError: handleError
})

const update = (user: User) => setTicketOwnerMutation.mutate(user)
</script>

<template>
  <div v-if="isPending" class="flex justify-center">
    <LoaderCircle class="h-16 w-16 animate-spin text-primary" />
  </div>
  <Alert v-else-if="isError" variant="destructive" class="mb-4">
    <AlertTitle>Error</AlertTitle>
    <AlertDescription>{{ error }}</AlertDescription>
  </Alert>
  <UserSelect v-if="!user" @update:modelValue="update">
    <Button variant="outline" role="combobox">
      <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
      Unassigned
    </Button>
  </UserSelect>
  <UserSelect v-else :modelValue="user" @update:modelValue="update">
    <Button variant="outline" role="combobox">
      <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
      {{ user.name }}
    </Button>
  </UserSelect>
</template>
