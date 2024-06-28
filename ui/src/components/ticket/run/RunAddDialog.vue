<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { toast } from '@/components/ui/toast'

import { Plus } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import type { Playbook, Run, Ticket } from '@/lib/types'

const queryClient = useQueryClient()

const props = defineProps<{
  ticket: Ticket
}>()

const newPlaybook = ref('')
const isOpen = ref(false)

const {
  isPending,
  isError,
  data: playbooks,
  error
} = useQuery({
  queryKey: ['playbooks'],
  queryFn: (): Promise<Array<Playbook>> =>
    pb.collection('playbooks').getFullList({
      sort: '-created'
    })
})

const addRunMutation = useMutation({
  mutationFn: (): Promise<Run> => {
    if (!playbooks.value) return Promise.reject('Playbooks not loaded')

    const playbook = playbooks.value.find((pb) => pb.id === newPlaybook.value)

    if (!playbook) return Promise.reject('Playbook not found')

    const steps = playbook.steps.map((step) => {
      return {
        name: step.name,
        type: step.type,
        status: 'open',
        description: step.description,
        schema: step.schema,
        state: {}
      }
    })

    return pb.collection('runs').create({
      ticket: props.ticket.id,
      name: playbook.name,
      steps: steps
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', props.ticket.id] })
    isOpen.value = false
    newPlaybook.value = ''
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

const submitDisabled = computed(() => isPending.value || isError.value || !newPlaybook.value)

const cancel = () => {
  newPlaybook.value = ''
  isOpen.value = false
}
</script>

<template>
  <Button v-if="!isOpen" variant="outline" @click="isOpen = true">
    <Plus class="mr-2 size-4" />
    Add Playbook
  </Button>
  <Card v-else class="flex flex-row items-center gap-2 p-2">
    <TanView :isError="isError" :isPending="isPending" :error="error" :value="playbooks">
      <Select v-model="newPlaybook">
        <SelectTrigger class="w-full">
          <SelectValue placeholder="Select Playbook" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem v-for="playbook in playbooks" :key="playbook.id" :value="playbook.id">
              {{ playbook.name }}
            </SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
    </TanView>
    <Button variant="outline" @click="cancel"> Cancel</Button>
    <Button :disabled="submitDisabled" @click="addRunMutation.mutate"> Add Playbook </Button>
  </Card>
</template>
