<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useToast } from '@/components/ui/toast/use-toast'
import UserForm from '@/components/user/UserForm.vue'
import UserGroup from '@/components/user/UserGroup.vue'
import UserPasswordForm from '@/components/user/UserPasswordForm.vue'

import { ChevronLeft } from 'lucide-vue-next'

import CardContent from '../ui/card/CardContent.vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { User, UserUpdate } from '@/client/models'
import { handleError } from '@/lib/utils'

// Prevent unused var warnings for components used in the template
const _tabsComponents = { Tabs, TabsContent, TabsList, TabsTrigger }
void _tabsComponents

const api = useAPI()

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const passwordForm = ref<InstanceType<typeof UserPasswordForm> | null>(null)

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
  onSuccess: () => {
    toast({
      title: 'User updated',
      description: 'The user has been updated successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['users'] })
  },
  onError: handleError('Failed to update user')
})

const passwordSubmit = (values: UserUpdate) => {
  updateUserMutation.mutate(values, {
    onSuccess: () => {
      passwordForm.value?.reset()
    }
  })
}

const deleteMutation = useMutation({
  mutationFn: () => api.deleteUser({ id: props.id }),
  onSuccess: () => {
    queryClient.removeQueries({ queryKey: ['users', props.id] })
    queryClient.invalidateQueries({ queryKey: ['users'] })
    toast({
      title: 'User deleted',
      description: 'The user has been deleted successfully'
    })
    router.push({ name: 'users' })
  },
  onError: handleError('Failed to delete user')
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
      <ColumnBodyContainer>
        <div class="flex flex-col gap-4 xl:flex-row">
          <div class="flex flex-col gap-4 xl:flex-1">
            <Card>
              <CardHeader>
                <CardTitle>User</CardTitle>
              </CardHeader>
              <CardContent>
                <UserForm :user="user" @submit="updateUserMutation.mutate" />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Set Password</CardTitle>
              </CardHeader>
              <CardContent>
                <UserPasswordForm ref="passwordForm" @submit="passwordSubmit" />
              </CardContent>
            </Card>
          </div>

          <div class="flex w-full flex-col gap-4 xl:w-96 xl:shrink-0">
            <Card>
              <CardHeader>
                <CardTitle>Access</CardTitle>
              </CardHeader>
              <CardContent>
                <UserGroup :id="user.id" />
              </CardContent>
            </Card>
          </div>
        </div>
        <!--Tabs default-value="groups" class="w-full">
          <TabsList class="grid w-full grid-cols-3">
            <TabsTrigger value="user"> User </TabsTrigger>
            <TabsTrigger value="password"> Password </TabsTrigger>
            <TabsTrigger value="groups"> Groups </TabsTrigger>
          </TabsList>
          <TabsContent value="user" class="mt-2">
            <Card>
              <CardHeader>
                <CardTitle>User</CardTitle>
              </CardHeader>
              <CardContent>
                <UserForm :user="user" @submit="updateUserMutation.mutate" />
              </CardContent>
            </Card>
          </TabsContent>
          <TabsContent value="password" class="mt-2">
            <Card>
              <CardHeader>
                <CardTitle>Set Password</CardTitle>
              </CardHeader>
              <CardContent>
                <UserPasswordForm ref="passwordForm" @submit="passwordSubmit" />
              </CardContent>
            </Card>
          </TabsContent>
          <TabsContent value="groups" class="mt-2">
            <UserGroup :id="user.id" />
          </TabsContent>
        </Tabs-->
      </ColumnBodyContainer>
    </ColumnBody>
  </TanView>
</template>
