<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import TypeForm from '@/components/type/TypeForm.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Type, TypeUpdate } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const { toast } = useToast()
const router = useRouter()
const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const {
  isPending,
  isError,
  data: type,
  error
} = useQuery({
  queryKey: ['types', props.id],
  queryFn: (): Promise<Type> => api.getType({ id: props.id })
})

const updateTypeMutation = useMutation({
  mutationFn: (update: TypeUpdate) => api.updateType({ id: props.id, typeUpdate: update }),
  onSuccess: () => {
    toast({
      title: 'Type updated',
      description: 'The type has been updated successfully'
    })
    queryClient.invalidateQueries({ queryKey: ['types'] })
  },
  onError: handleError('Failed to update type')
})

const deleteMutation = useMutation({
  mutationFn: () => api.deleteType({ id: props.id }),
  onSuccess: () => {
    queryClient.removeQueries({ queryKey: ['types', props.id] })
    queryClient.invalidateQueries({ queryKey: ['types'] })
    queryClient.invalidateQueries({ queryKey: ['sidebar'] })
    toast({
      title: 'Type deleted',
      description: 'The type has been deleted successfully'
    })
    router.push({ name: 'types' })
  },
  onError: handleError('Failed to delete type')
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <ColumnHeader>
      <Button @click="router.push({ name: 'types' })" variant="outline" class="sm:hidden">
        <ChevronLeft class="mr-2 size-4" />
        Back
      </Button>
      <div class="ml-auto">
        <DeleteDialog
          v-if="type"
          :name="type.singular"
          singular="Type"
          @delete="deleteMutation.mutate"
        />
      </div>
    </ColumnHeader>

    <ColumnBody v-if="type">
      <ColumnBodyContainer small>
        <TypeForm :type="type" @submit="updateTypeMutation.mutate" />
      </ColumnBodyContainer>
    </ColumnBody>
  </TanView>
</template>
