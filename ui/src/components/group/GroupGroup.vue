<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import GroupSelectDialog from '@/components/group/GroupSelectDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useToast } from '@/components/ui/toast/use-toast'
import UserSelectDialog from '@/components/user/UserSelectDialog.vue'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { GroupUser, UserGroup } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()
const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  id: string
}>()

const { data: parentGroups } = useQuery({
  queryKey: ['parent_groups', props.id],
  queryFn: (): Promise<Array<UserGroup>> => api.listParentGroups({ id: props.id })
})

const { data: parentPermissions } = useQuery({
  queryKey: ['parent_permissions', props.id],
  queryFn: (): Promise<Array<string>> => api.listParentPermissions({ id: props.id })
})

const { data: childGroups } = useQuery({
  queryKey: ['child_groups', props.id],
  queryFn: (): Promise<Array<UserGroup>> => api.listChildGroups({ id: props.id })
})

const { data: groupUsers } = useQuery({
  queryKey: ['group_users', props.id],
  queryFn: (): Promise<Array<GroupUser>> => api.listGroupUsers({ id: props.id })
})

const addGroupUserMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addUserGroup({
      id: id,
      groupRelation: {
        groupId: props.id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['group_users'] })
    toast({
      title: 'User added',
      description: 'The user has been added to the group'
    })
  },
  onError: handleError('Failed to add user to group')
})

const addGroupParentMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addGroupParent({
      id: id,
      groupRelation: {
        groupId: props.id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['parent_groups'] })
    queryClient.invalidateQueries({ queryKey: ['parent_permissions'] })
    toast({
      title: 'Parent group added',
      description: 'The parent group has been added successfully'
    })
  },
  onError: handleError('Failed to add parent group')
})

const addGroupChildMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addGroupParent({
      id: props.id,
      groupRelation: {
        groupId: id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['child_groups'] })
    toast({
      title: 'Child group added',
      description: 'The child group has been added successfully'
    })
  },
  onError: handleError('Failed to add child group')
})

const removeGroupUserMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeUserGroup({
      id: id,
      groupId: props.id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['group_users'] })
    toast({
      title: 'User removed',
      description: 'The user has been removed from the group'
    })
  },
  onError: handleError('Failed to remove user from group')
})

const removeGroupParentMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeGroupParent({
      id: id,
      parentGroupId: props.id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['parent_groups'] })
    queryClient.invalidateQueries({ queryKey: ['parent_permissions'] })
    toast({
      title: 'Parent group removed',
      description: 'The parent group has been removed successfully'
    })
  },
  onError: handleError('Failed to remove parent group')
})

const removeGroupChildMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeGroupParent({
      id: props.id,
      parentGroupId: id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['child_groups'] })
    toast({
      title: 'Child group removed',
      description: 'The child group has been removed successfully'
    })
  },
  onError: handleError('Failed to remove child group')
})

const dialogOpenParent = ref(false)
const dialogOpenChild = ref(false)
const dialogOpenUser = ref(false)

const selectParent = (group: { group: string }) => {
  addGroupParentMutation.mutate(group.group)
  dialogOpenParent.value = false
}

const selectChild = (group: { group: string }) => {
  addGroupChildMutation.mutate(group.group)
  dialogOpenChild.value = false
}

const selectUser = (user: { user: string }) => {
  addGroupUserMutation.mutate(user.user)
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
        <TicketPanel title="Child Groups" @add="dialogOpenChild = true">
          <GroupSelectDialog
            v-model="dialogOpenChild"
            @select="selectChild"
            :exclude="childGroups?.map((group) => group.id).concat([id]) ?? [id]"
          />
          <PanelListElement
            v-for="groupGroup in childGroups"
            :key="groupGroup.id"
            class="flex h-10 flex-row items-center pr-1"
          >
            <div class="flex flex-1 items-center overflow-hidden">
              <RouterLink
                :to="{ name: 'groups', params: { id: groupGroup.id } }"
                class="hover:underline"
              >
                {{ groupGroup.name }}
              </RouterLink>
              <span class="ml-1 text-sm text-muted-foreground">({{ groupGroup.type }})</span>
            </div>
            <DeleteDialog
              v-if="groupGroup.type === 'direct'"
              :name="groupGroup.name"
              singular="Membership"
              @delete="removeGroupChildMutation.mutate(groupGroup.id)"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
          <div
            v-if="!childGroups || childGroups.length === 0"
            class="flex h-10 items-center p-4 text-muted-foreground"
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
            :exclude="groupUsers?.map((user) => user.id) ?? []"
          />
          <PanelListElement
            v-for="groupUser in groupUsers"
            :key="groupUser.id"
            class="flex h-10 flex-row items-center pr-1"
          >
            <div class="flex flex-1 items-center overflow-hidden">
              <RouterLink
                :to="{ name: 'users', params: { id: groupUser.id } }"
                class="hover:underline"
              >
                {{ groupUser.name }}
              </RouterLink>
              <span class="ml-1 text-sm text-muted-foreground">({{ groupUser.type }})</span>
            </div>
            <DeleteDialog
              v-if="groupUser.type === 'direct'"
              :name="groupUser.name"
              singular="Membership"
              @delete="removeGroupUserMutation.mutate(groupUser.id)"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
          <div
            v-if="!groupUsers || groupUsers.length === 0"
            class="flex h-10 items-center p-4 text-muted-foreground"
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
        <TicketPanel title="Parent Groups" @add="dialogOpenParent = true">
          <GroupSelectDialog
            v-model="dialogOpenParent"
            @select="selectParent"
            :exclude="parentGroups?.map((group) => group.id).concat([id]) ?? [id]"
          />
          <PanelListElement
            v-for="groupGroup in parentGroups"
            :key="groupGroup.id"
            class="flex h-10 flex-row items-center pr-1"
          >
            <div class="flex flex-1 items-center overflow-hidden">
              <RouterLink
                :to="{ name: 'groups', params: { id: groupGroup.id } }"
                class="hover:underline"
              >
                {{ groupGroup.name }}
              </RouterLink>
              <span class="ml-1 text-sm text-muted-foreground">({{ groupGroup.type }})</span>
            </div>
            <DeleteDialog
              v-if="groupGroup.type === 'direct'"
              :name="groupGroup.name"
              singular="Inheritance"
              @delete="removeGroupParentMutation.mutate(groupGroup.id)"
            >
              <Button variant="ghost" size="icon" class="h-8 w-8">
                <Trash2 class="size-4" />
              </Button>
            </DeleteDialog>
          </PanelListElement>
          <div
            v-if="!parentGroups || parentGroups.length === 0"
            class="flex h-10 items-center p-4 text-muted-foreground"
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
