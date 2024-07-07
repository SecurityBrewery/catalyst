<script setup lang="ts">
import UserSelect from '@/components/common/UserSelect.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { toast } from '@/components/ui/toast'

import { LoaderCircle, User2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Ticket, User } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
  uID: string
}>()

const {
  isPending,
  isError,
  data: user,
  error
} = useQuery({
  queryKey: ['tickets', props.ticket.id, 'owner', props.uID],
  queryFn: (): Promise<User> => pb.collection('users').getOne(props.uID)
})

const setTicketOwnerMutation = useMutation({
  mutationFn: (user: User): Promise<Ticket> =>
    pb.collection('tickets').update(props.ticket.id, {
      owner: user.id
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets'] }),
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
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
  <div v-if="!user">
    <Button variant="outline" role="combobox" disabled>
      <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
      {{ props.uID }}
    </Button>
  </div>
  <UserSelect v-else :modelValue="user" @update:modelValue="update">
    <Button variant="outline" role="combobox">
      <User2 class="mr-2 size-4 h-4 w-4 shrink-0 opacity-50" />
      {{ user.name }}
    </Button>
  </UserSelect>
</template>