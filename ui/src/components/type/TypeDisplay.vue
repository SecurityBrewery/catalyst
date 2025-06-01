<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import { Button } from '@/components/ui/button'
import TypeForm from '@/components/type/TypeForm.vue'

import { ChevronLeft } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { Type, TypeUpdate } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

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
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['types'] }),
  onError: handleError
})

const deleteMutation = useMutation({
  mutationFn: () => {
    return api.deleteType({ id: props.id })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['types'] })
    router.push({ name: 'types' })
  },
  onError: handleError
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
          :name="type.name"
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
