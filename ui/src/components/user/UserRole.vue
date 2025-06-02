<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import RoleSelectDialog from '@/components/role/RoleSelectDialog.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

import { useAPI } from '@/api'
import type { Role, UserRole } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()
const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const { data: userRoles } = useQuery({
  queryKey: ['user_roles', props.id],
  queryFn: (): Promise<Array<UserRole>> => api.listUserRoles({ id: props.id })
})

const { data: userPermissions } = useQuery({
  queryKey: ['user_permissions', props.id],
  queryFn: (): Promise<Array<string>> => api.listUserPermissions({ id: props.id })
})

const addRoleMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addUserRole({
      id: props.id,
      roleRelation: {
        roleId: id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['user_roles'] })
    queryClient.invalidateQueries({ queryKey: ['user_permissions'] })
  },
  onError: handleError
})

const removeRoleMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeUserRole({
      id: props.id,
      roleId: id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['user_roles'] })
    queryClient.invalidateQueries({ queryKey: ['user_permissions'] })
  },
  onError: handleError
})

const dialogOpen = ref(false)

const select = (role: Role) => {
  addRoleMutation.mutate(role.id)
  dialogOpen.value = false
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <TicketPanel title="Assigned Roles" @add="dialogOpen = true">
      <RoleSelectDialog
        v-model="dialogOpen"
        @select="select"
        :exclude="userRoles?.map((role) => role.id) ?? []"
      />
      <PanelListElement
        v-for="userRole in userRoles"
        :key="userRole.id"
        class="flex-row items-center pr-1"
      >
        <div class="flex flex-1 items-center overflow-hidden">
          {{ userRole.name }}
          <span class="ml-1 text-sm text-muted-foreground">({{ userRole.type }})</span>
        </div>
        <DeleteDialog
          :name="userRole.name"
          singular="Role"
          @delete="removeRoleMutation.mutate(userRole.id)"
        >
          <Button variant="ghost" size="icon" class="h-8 w-8">
            <Trash2 class="size-4" />
          </Button>
        </DeleteDialog>
      </PanelListElement>
      <div
        v-if="!userRoles || userRoles.length === 0"
        class="flex h-10 items-center p-4 text-sm text-muted-foreground"
      >
        No roles assigned yet.
      </div>
    </TicketPanel>
  </div>

  <div class="mt-4 flex flex-col gap-4">
    <h2 class="text-sm font-medium">User Permissions</h2>
    <p class="text-sm text-muted-foreground">
      The following permissions are granted to the user by their roles.
    </p>
    <div class="flex flex-wrap gap-2">
      <Badge v-for="(permission, index) in userPermissions" :key="index">{{ permission }}</Badge>
    </div>
  </div>
</template>
