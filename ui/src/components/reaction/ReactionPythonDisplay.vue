<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import DeleteDialog from '@/components/common/DeleteDialog.vue'
import DynamicInput from '@/components/input/DynamicInput.vue'
import ReactionPythonForm from '@/components/reaction/ReactionPythonForm.vue'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { toast } from '@/components/ui/toast'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { ReactionPython } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  id: string
}>()

const {
  isPending,
  isError,
  data: reaction,
  error
} = useQuery({
  queryKey: ['reactions_python', props.id],
  queryFn: (): Promise<ReactionPython> => pb.collection('reactions_python').getOne(props.id)
})

const updateReactionMutation = useMutation({
  mutationFn: (update: any) => pb.collection('reactions_python').update(props.id, update),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    queryClient.invalidateQueries({ queryKey: ['reactions_python'] })
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error" :value="reaction">
    <div class="flex h-full flex-1 flex-col overflow-hidden">
      <div class="flex items-center bg-background px-4 py-2">
        <div class="ml-auto">
          <DeleteDialog
            v-if="reaction"
            collection="reactions_python"
            :id="reaction.id"
            :name="reaction.name"
            :singular="'Reaction'"
            :to="{ name: 'reactions' }"
            :queryKey="['reactions']"
          />
        </div>
      </div>
      <Separator />

      <ScrollArea v-if="reaction" class="flex-1">
        <div class="flex max-w-[640px] flex-col gap-4 p-4">
          <h1 class="text-3xl font-bold">Python Reaction Editor</h1>

          <p class="py-4 text-sm text-muted-foreground">
            Write a Python script that will be executed when this reaction is triggered.
          </p>

          <ReactionPythonForm :reaction="reaction" @submit="updateReactionMutation.mutate" />
        </div>
      </ScrollArea>
    </div>
  </TanView>
</template>
