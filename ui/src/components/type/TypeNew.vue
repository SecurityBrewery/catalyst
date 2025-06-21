<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import TypeForm from '@/components/type/TypeForm.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { NewType, Type } from '@/client'
import { handleError } from '@/lib/utils'

const api = useAPI()

const queryClient = useQueryClient()
const router = useRouter()
const { toast } = useToast()

const addTypeMutation = useMutation({
  mutationFn: (values: NewType): Promise<Type> => api.createType({ newType: values }),
  onSuccess: (data: Type) => {
    router.push({ name: 'types', params: { id: data.id } })
    toast({
      title: 'Type created',
      description: 'The type has been created successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['types'] })
  },
  onError: handleError('Failed to create type')
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
