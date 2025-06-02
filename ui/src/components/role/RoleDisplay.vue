<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import RoleForm from '@/components/role/RoleForm.vue'
import RoleRole from '@/components/role/RoleRole.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Role, RoleUpdate } from '@/client/models'
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
  data: role,
  error
} = useQuery({
  queryKey: ['roles', props.id],
  queryFn: (): Promise<Role> => api.getRole({ id: props.id })
})

const updateRoleMutation = useMutation({
  mutationFn: (update: RoleUpdate) => api.updateRole({ id: props.id, roleUpdate: update }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['roles'] }),
  onError: handleError
})

const deleteMutation = useMutation({
  mutationFn: () => {
    return api.deleteRole({ id: props.id })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['roles'] })
    router.push({ name: 'groups' })
  },
  onError: handleError
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader>
      <Button @click="router.push({ name: 'groups' })" variant="outline" class="sm:hidden">
        <ChevronLeft class="mr-2 size-4" />
        Back
      </Button>
      <div class="ml-auto">
        <DeleteDialog
          v-if="role"
          :name="role.name"
          singular="Group"
          @delete="deleteMutation.mutate"
        />
      </div>
    </ColumnHeader>

    <ColumnBody v-if="role">
      <ColumnBodyContainer>
        <div class="flex flex-col gap-4 xl:flex-row">
          <div class="flex flex-col gap-4 xl:flex-1">
            <RoleForm :role="role" @submit="updateRoleMutation.mutate" />
          </div>
          <div class="flex w-full flex-col gap-4 xl:w-96 xl:shrink-0">
            <RoleRole :id="role.id" />
          </div>
        </div>
      </ColumnBodyContainer>
    </ColumnBody>
  </TanView>
</template>
