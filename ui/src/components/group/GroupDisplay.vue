<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import GroupForm from '@/components/group/GroupForm.vue'
import GroupGroup from '@/components/group/GroupGroup.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Group, GroupUpdate } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const props = defineProps<{
  id: string
}>()

const {
  isPending,
  isError,
  data: group,
  error
} = useQuery({
  queryKey: ['groups', props.id],
  queryFn: (): Promise<Group> => api.getGroup({ id: props.id })
})

const updateGroupMutation = useMutation({
  mutationFn: (update: GroupUpdate) => api.updateGroup({ id: props.id, groupUpdate: update }),
  onSuccess: () => {
    toast({
      title: 'Group updated',
      description: 'The group has been updated successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['groups'] })
  },
  onError: handleError('Failed to update group')
})

const deleteMutation = useMutation({
  mutationFn: () => api.deleteGroup({ id: props.id }),
  onSuccess: () => {
    queryClient.removeQueries({ queryKey: ['groups', props.id] })
    queryClient.invalidateQueries({ queryKey: ['groups'] })
    toast({
      title: 'Group deleted',
      description: 'The group has been deleted successfully'
    })
    router.push({ name: 'groups' })
  },
  onError: handleError('Failed to delete group')
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
          v-if="group && group.id !== 'admin'"
          :name="group.name"
          singular="Group"
          @delete="deleteMutation.mutate"
        />
      </div>
    </ColumnHeader>

    <ColumnBody v-if="group">
      <ColumnBodyContainer>
        <div class="flex flex-col gap-4 xl:flex-row">
          <div class="flex flex-col gap-4 xl:flex-1">
            <GroupForm :group="group" @submit="updateGroupMutation.mutate" />
          </div>
          <div class="flex w-full flex-col gap-4 xl:w-96 xl:shrink-0">
            <GroupGroup :id="group.id" />
          </div>
        </div>
      </ColumnBodyContainer>
    </ColumnBody>
  </TanView>
</template>
