<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import UserForm from '@/components/user/UserForm.vue'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { User, UserUpdate } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const router = useRouter()
const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const {
  isPending,
  isError,
  data: user,
  error
} = useQuery({
  queryKey: ['users', props.id],
  queryFn: (): Promise<User> => api.getUser({ id: props.id })
})

const updateUserMutation = useMutation({
  mutationFn: (update: UserUpdate) => api.updateUser({ id: props.id, userUpdate: update }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['users'] }),
  onError: handleError
})

const deleteMutation = useMutation({
  mutationFn: () => {
    return api.deleteUser({ id: props.id })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['users'] })
    router.push({ name: 'users' })
  },
  onError: handleError
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader>
      <Button @click="router.push({ name: 'users' })" variant="outline" class="sm:hidden">
        <ChevronLeft class="mr-2 size-4" />
        Back
      </Button>
      <div class="ml-auto">
        <DeleteDialog
          v-if="user"
          :name="user.name"
          singular="User"
          @delete="deleteMutation.mutate"
        />
      </div>
    </ColumnHeader>

    <ColumnBody v-if="user">
      <ColumnBodyContainer small>
        <UserForm :user="user" @submit="updateUserMutation.mutate" />
      </ColumnBodyContainer>
    </ColumnBody>
  </TanView>
</template>
