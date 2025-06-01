<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import TypeForm from '@/components/type/TypeForm.vue'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { NewType, Type } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()

const addTypeMutation = useMutation({
  mutationFn: (values: NewType): Promise<Type> => api.createType({ newType: values }),
  onSuccess: (data: Type) => {
    router.push({ name: 'types', params: { id: data.id } })
    queryClient.invalidateQueries({ queryKey: ['types'] })
  },
  onError: handleError
})
</script>

<template>
  <ColumnHeader>
    <Button @click="router.push({ name: 'types' })" variant="outline" class="sm:hidden">
      <ChevronLeft class="mr-2 size-4" />
      Back
    </Button>
  </ColumnHeader>

  <ColumnBody>
    <ColumnBodyContainer small>
      <TypeForm @submit="addTypeMutation.mutate" />
    </ColumnBodyContainer>
  </ColumnBody>
</template>
