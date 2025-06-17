<script setup lang="ts">
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import GroupSelectDialog from '@/components/group/GroupSelectDialog.vue'
import PanelListElement from '@/components/layout/PanelListElement.vue'
import TicketPanel from '@/components/ticket/TicketPanel.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { Trash2 } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'

import { useAPI } from '@/api'
import type { UserGroup } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()
const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  id: string
}>()

const { data: userGroups } = useQuery({
  queryKey: ['user_groups', props.id],
  queryFn: (): Promise<Array<UserGroup>> => api.listUserGroups({ id: props.id })
})

const { data: userPermissions } = useQuery({
  queryKey: ['user_permissions', props.id],
  queryFn: (): Promise<Array<string>> => api.listUserPermissions({ id: props.id })
})

const addGroupMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.addUserGroup({
      id: props.id,
      groupRelation: {
        groupId: id
      }
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['user_groups'] })
    queryClient.invalidateQueries({ queryKey: ['user_permissions'] })
    toast({
      title: 'Group added',
      description: 'The group has been added successfully'
    })
  },
  onError: handleError
})

const removeGroupMutation = useMutation({
  mutationFn: (id: string): Promise<void> =>
    api.removeUserGroup({
      id: props.id,
      groupId: id
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['user_groups'] })
    queryClient.invalidateQueries({ queryKey: ['user_permissions'] })
    toast({
      title: 'Group removed',
      description: 'The group has been removed successfully'
    })
  },
  onError: handleError
})

const dialogOpen = ref(false)

const select = (group: { group: string }) => {
  addGroupMutation.mutate(group.group)
  dialogOpen.value = false
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <TicketPanel title="Groups" @add="dialogOpen = true">
      <GroupSelectDialog
        v-model="dialogOpen"
        @select="select"
        :exclude="userGroups?.map((group) => group.id) ?? []"
      />
      <PanelListElement
        v-for="userGroup in userGroups"
        :key="userGroup.id"
        class="flex h-10 flex-row items-center pr-1"
      >
        <div class="flex flex-1 items-center overflow-hidden">
          <RouterLink
            :to="{ name: 'groups', params: { id: userGroup.id } }"
            class="hover:underline"
          >
            {{ userGroup.name }}
          </RouterLink>
          <span class="ml-1 text-sm text-muted-foreground">({{ userGroup.type }})</span>
        </div>
        <DeleteDialog
          v-if="userGroup.type === 'direct'"
          :name="userGroup.name"
          singular="Group Membership"
          @delete="removeGroupMutation.mutate(userGroup.id)"
        >
          <Button variant="ghost" size="icon" class="h-8 w-8">
            <Trash2 class="size-4" />
          </Button>
        </DeleteDialog>
      </PanelListElement>
      <div
        v-if="!userGroups || userGroups.length === 0"
        class="flex h-10 items-center p-4 text-muted-foreground"
      >
        No groups assigned yet.
      </div>
    </TicketPanel>
  </div>

  <div class="mt-4 flex flex-col gap-4">
    <h2 class="text-sm font-medium">Permissions</h2>
    <p class="text-sm text-muted-foreground">
      The following permissions are granted to the user by their groups.
    </p>
    <div class="flex flex-wrap gap-2">
      <Badge v-for="(permission, index) in userPermissions" :key="index">{{ permission }}</Badge>
    </div>
  </div>
</template>
