<script setup lang="ts">
import GroupForm from '@/components/group/GroupForm.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Group, NewGroup } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()
const { toast } = useToast()

const addGroupMutation = useMutation({
  mutationFn: (values: NewGroup): Promise<Group> => api.createGroup({ newGroup: values }),
  onSuccess: (data: Group) => {
    router.push({ name: 'groups', params: { id: data.id } })
    toast({
      title: 'Group created',
      description: 'The group has been created successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['groups'] })
  },
  onError: handleError('Failed to create group')
})
</script>

<template>
  <ColumnHeader>
    <Button @click="router.push({ name: 'groups' })" variant="outline" class="sm:hidden">
      <ChevronLeft class="mr-2 size-4" />
      Back
    </Button>
  </ColumnHeader>

  <ColumnBody>
    <ColumnBodyContainer small>
      <GroupForm @submit="addGroupMutation.mutate" />
    </ColumnBodyContainer>
  </ColumnBody>
</template>
