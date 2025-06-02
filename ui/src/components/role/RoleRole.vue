<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import RoleSelectDialog from '@/components/role/RoleSelectDialog.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import UserSelectDialog from '@/components/user/UserSelectDialog.vue'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { RoleUser, UserRole } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()
const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const { data: parentRoles } = useQuery({
  queryKey: ['parent_roles', props.id],
  queryFn: (): Promise<Array<UserRole>> => api.listParentRoles({ id: props.id })
})

const { data: parentPermissions } = useQuery({
  queryKey: ['parent_permissions', props.id],
  queryFn: (): Promise<Array<string>> => api.listParentPermissions({ id: props.id })
})

const { data: childRoles } = useQuery({
  queryKey: ['child_roles', props.id],
  queryFn: (): Promise<Array<UserRole>> => api.listChildRoles({ id: props.id })
})

const { data: roleUsers } = useQuery({
  queryKey: ['role_users', props.id],
  queryFn: (): Promise<Array<RoleUser>> => api.listRoleUsers({ id: props.id })
})

const addRoleUserMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addUserRole({
      id: id,
      roleRelation: {
        roleId: props.id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['role_users'] })
  },
  onError: handleError
})

const addRoleParentMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addRoleParent({
      id: id,
      roleRelation: {
        roleId: props.id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['parent_roles'] })
    queryClient.invalidateQueries({ queryKey: ['parent_permissions'] })
  },
  onError: handleError
})

const addRoleChildMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addRoleParent({
      id: props.id,
      roleRelation: {
        roleId: id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['child_roles'] })
  },
  onError: handleError
})

const removeRoleUserMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeUserRole({
      id: id,
      roleId: props.id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['role_users'] })
  },
  onError: handleError
})

const removeRoleParentMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeRoleParent({
      id: id,
      parentRoleId: props.id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['parent_roles'] })
    queryClient.invalidateQueries({ queryKey: ['parent_permissions'] })
  },
  onError: handleError
})

const removeRoleChildMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeRoleParent({
      id: props.id,
      parentRoleId: id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['child_roles'] })
  },
  onError: handleError
})

const dialogOpenParent = ref(false)
const dialogOpenChild = ref(false)
const dialogOpenUser = ref(false)

const selectParent = (role: { role: string }) => {
  addRoleParentMutation.mutate(role.role)
  dialogOpenParent.value = false
}

const selectChild = (role: { role: string }) => {
  addRoleChildMutation.mutate(role.role)
  dialogOpenChild.value = false
}

const selectUser = (user: { user: string }) => {
  addRoleUserMutation.mutate(user.user)
  dialogOpenUser.value = false
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Members</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="flex flex-col gap-4">
        <TicketPanel title="Groups" @add="dialogOpenChild = true">
          <RoleSelectDialog
            v-model="dialogOpenChild"
            @select="selectChild"
            :exclude="childRoles?.map((role) => role.id).concat([id]) ?? [id]"
          />
          <PanelListElement
            v-for="roleRole in childRoles"
            :key="roleRole.id"
            class="flex h-10 flex-row items-center pr-1"
          >
            <div class="flex flex-1 items-center overflow-hidden">
              <RouterLink
                :to="{ name: 'groups', params: { id: roleRole.id } }"
                class="hover:underline"
              >
                {{ roleRole.name }}
              </RouterLink>
              <span class="ml-1 text-sm text-muted-foreground">({{ roleRole.type }})</span>
            </div>
            <DeleteDialog
              v-if="roleRole.type === 'direct'"
              :name="roleRole.name"
              singular="Group"
              @delete="removeRoleChildMutation.mutate(roleRole.id)"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
          <div
            v-if="!childRoles || childRoles.length === 0"
            class="flex h-10 items-center p-4 text-sm text-muted-foreground"
          >
            No groups assigned yet.
          </div>
        </TicketPanel>
      </div>

      <div class="mt-4 flex flex-col gap-4">
        <TicketPanel title="Users" @add="dialogOpenUser = true">
          <UserSelectDialog
            v-model="dialogOpenUser"
            @select="selectUser"
            :exclude="roleUsers?.map((user) => user.id) ?? []"
          />
          <PanelListElement
            v-for="roleUser in roleUsers"
            :key="roleUser.id"
            class="flex h-10 flex-row items-center pr-1"
          >
            <div class="flex flex-1 items-center overflow-hidden">
              <RouterLink
                :to="{ name: 'users', params: { id: roleUser.id } }"
                class="hover:underline"
              >
                {{ roleUser.name }}
              </RouterLink>
              <span class="ml-1 text-sm text-muted-foreground">({{ roleUser.type }})</span>
            </div>
            <DeleteDialog
              v-if="roleUser.type === 'direct'"
              :name="roleUser.name"
              singular="User"
              @delete="removeRoleUserMutation.mutate(roleUser.id)"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
          <div
            v-if="!roleUsers || roleUsers.length === 0"
            class="flex h-10 items-center p-4 text-sm text-muted-foreground"
          >
            No users assigned yet.
          </div>
        </TicketPanel>
      </div>
    </CardContent>
  </Card>

  <Card>
    <CardHeader>
      <CardTitle>Inheritance</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="mt-4 flex flex-col gap-4">
        <TicketPanel title="Groups" @add="dialogOpenParent = true">
          <RoleSelectDialog
            v-model="dialogOpenParent"
            @select="selectParent"
            :exclude="parentRoles?.map((role) => role.id).concat([id]) ?? [id]"
          />
          <PanelListElement
            v-for="roleRole in parentRoles"
            :key="roleRole.id"
            class="flex h-10 flex-row items-center pr-1"
          >
            <div class="flex flex-1 items-center overflow-hidden">
              <RouterLink
                :to="{ name: 'groups', params: { id: roleRole.id } }"
                class="hover:underline"
              >
                {{ roleRole.name }}
              </RouterLink>
              <span class="ml-1 text-sm text-muted-foreground">({{ roleRole.type }})</span>
            </div>
            <DeleteDialog
              v-if="roleRole.type === 'direct'"
              :name="roleRole.name"
              singular="Group"
              @delete="removeRoleParentMutation.mutate(roleRole.id)"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
          <div
            v-if="!parentRoles || parentRoles.length === 0"
            class="flex h-10 items-center p-4 text-sm text-muted-foreground"
          >
            No groups assigned yet.
          </div>
        </TicketPanel>
      </div>

      <div class="mt-4 flex flex-col gap-4">
        <h2 class="text-sm font-medium">Permissions</h2>
        <p class="text-sm text-muted-foreground">
          The following permissions are granted in addition to the permissions selected to the left
          by the parent groups.
        </p>
        <div class="flex flex-wrap gap-2">
          <Badge v-for="(permission, index) in parentPermissions" :key="index">{{
            permission
          }}</Badge>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
