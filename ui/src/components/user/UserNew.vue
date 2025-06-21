<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'
import UserForm from '@/components/user/UserForm.vue'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { NewUser, User } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()
const { toast } = useToast()

const addUserMutation = useMutation({
  mutationFn: (values: NewUser): Promise<User> => api.createUser({ newUser: values }),
  onSuccess: (data: User) => {
    router.push({ name: 'users', params: { id: data.id } })
    toast({
      title: 'User created',
      description: 'The user has been created successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['users'] })
  },
  onError: handleError('Failed to create user')
})
</script>

<template>
  <ColumnHeader>
    <Button @click="router.push({ name: 'users' })" variant="outline" class="sm:hidden">
      <ChevronLeft class="mr-2 size-4" />
      Back
    </Button>
  </ColumnHeader>

  <ColumnBody>
    <ColumnBodyContainer small>
      <UserForm @submit="addUserMutation.mutate" />
    </ColumnBodyContainer>
  </ColumnBody>
</template>
