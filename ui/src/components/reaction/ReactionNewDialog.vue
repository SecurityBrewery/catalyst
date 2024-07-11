<script setup lang="ts">
import ReactionPythonForm from '@/components/reaction/ReactionPythonForm.vue'
import ReactionWebhookForm from '@/components/reaction/ReactionWebhookForm.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { toast } from '@/components/ui/toast'

import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { ReactionPython, ReactionWebhook, Ticket } from '@/lib/types'

const queryClient = useQueryClient()
const router = useRouter()

const isOpen = ref(false)

const addReactionWebhookMutation = useMutation({
  mutationFn: (values: ReactionWebhook): Promise<ReactionWebhook> =>
    pb.collection('reactions_webhooks').create(values),
  onSuccess: (data: Ticket) => {
    router.push({ name: 'reactions', params: { type: 'webhook', id: data.id } })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    queryClient.invalidateQueries({ queryKey: ['reactions_webhooks'] })
    isOpen.value = false
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

const addReactionPythonMutation = useMutation({
  mutationFn: (values: ReactionPython): Promise<ReactionPython> =>
    pb.collection('reactions_python').create(values),
  onSuccess: (data: Ticket) => {
    router.push({ name: 'reactions', params: { type: 'python', id: data.id } })
    queryClient.invalidateQueries({ queryKey: ['reactions'] })
    queryClient.invalidateQueries({ queryKey: ['reactions_python'] })
    isOpen.value = false
  },
  onError: (error) =>
    toast({
      title: error.name,
      description: error.message,
      variant: 'destructive'
    })
})

const submitDisabled = ref(true)

const type = ref('')
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="ghost">New Reaction</Button>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New Reaction</DialogTitle>
        <DialogDescription>Create a new reaction</DialogDescription>
      </DialogHeader>

      <Label>Type</Label>
      <Select v-model="type">
        <SelectTrigger class="font-medium">
          <SelectValue placeholder="Select a type" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem value="python">Python</SelectItem>
            <SelectItem value="webhook">Webhook</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>

      <ReactionPythonForm
        v-if="type === 'python'"
        @submit="addReactionPythonMutation.mutate"
        v-model:submitDisabled="submitDisabled"
      >
        <DialogFooter class="mt-2">
          <DialogClose as-child>
            <Button type="button" variant="secondary" @click="isOpen = false">Cancel</Button>
          </DialogClose>
          <Button
            :title="submitDisabled ? 'Please fill out all required fields' : undefined"
            :disabled="submitDisabled"
            type="submit"
            >Save
          </Button>
        </DialogFooter>
      </ReactionPythonForm>
      <ReactionWebhookForm
        v-if="type === 'webhook'"
        @submit="addReactionWebhookMutation.mutate"
        v-model:submitDisabled="submitDisabled"
      >
        <DialogFooter class="mt-2">
          <DialogClose as-child>
            <Button type="button" variant="secondary">Cancel</Button>
          </DialogClose>
          <Button
            :title="submitDisabled ? 'Please fill out all required fields' : undefined"
            :disabled="submitDisabled"
            type="submit"
            >Save
          </Button>
        </DialogFooter>
      </ReactionWebhookForm>
    </DialogContent>
  </Dialog>
</template>
