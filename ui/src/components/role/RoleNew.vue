<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import RoleForm from '@/components/role/RoleForm.vue'
import { Button } from '@/components/ui/button'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { NewRole, Role } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()

const addRoleMutation = useMutation({
  mutationFn: (values: NewRole): Promise<Role> => api.createRole({ newRole: values }),
  onSuccess: (data: Role) => {
    router.push({ name: 'roles', params: { id: data.id } })
    queryClient.invalidateQueries({ queryKey: ['roles'] })
  },
  onError: handleError
})
</script>

<template>
  <ColumnHeader>
    <Button @click="router.push({ name: 'roles' })" variant="outline" class="sm:hidden">
      <ChevronLeft class="mr-2 size-4" />
      Back
    </Button>
  </ColumnHeader>

  <ColumnBody>
    <ColumnBodyContainer small>
      <RoleForm @submit="addRoleMutation.mutate" />
    </ColumnBodyContainer>
  </ColumnBody>
</template>
