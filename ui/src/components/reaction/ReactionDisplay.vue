<script setup lang="ts">
import TanView from '@/components/TanView.vue'
import GrowTextarea from '@/components/form/GrowTextarea.vue'
import DynamicInput from '@/components/input/DynamicInput.vue'
import ReactionDeleteDialog from '@/components/reaction/ReactionDeleteDialog.vue'
import ReactionPython from '@/components/reaction/ReactionPython.vue'
import ReactionWebhook from '@/components/reaction/ReactionWebhook.vue'
import TriggerHook from '@/components/reaction/TriggerHook.vue'
import TriggerTicketType from '@/components/reaction/TriggerTicketType.vue'
import TriggerWebhook from '@/components/reaction/TriggerWebhook.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { pb } from '@/lib/pocketbase'
import type { Reaction } from '@/lib/types'
import { handleError } from '@/lib/utils'

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
  queryKey: ['reactions', props.id],
  queryFn: (): Promise<Reaction> => pb.collection('reactions').getOne(props.id)
})

const updateReactionMutation = useMutation({
  mutationFn: (update: any) => pb.collection('reactions').update(props.id, update),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['reactions'] }),
  onError: handleError
})

const updateName = (name: string) => updateReactionMutation.mutate({ name: name })
const updateDescription = (description: string) =>
  updateReactionMutation.mutate({ description: description })
const updateTrigger = (trigger: string) =>
  updateReactionMutation.mutate({ trigger, triggerdata: {} })
const updateTriggerData = (triggerdata: any) => updateReactionMutation.mutate({ triggerdata })
const updateReaction = (reaction: string) =>
  updateReactionMutation.mutate({ reaction, reactiondata: {} })
const updateReactionData = (reactiondata: any) => updateReactionMutation.mutate({ reactiondata })
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error" :value="reaction">
    <div class="flex h-full flex-1 flex-col overflow-hidden">
      <div class="flex items-center bg-background px-4 py-2">
        <div class="ml-auto">
          <ReactionDeleteDialog v-if="reaction" :reaction="reaction" />
        </div>
      </div>
      <Separator />

      <ScrollArea v-if="reaction" class="flex-1">
        <div class="flex max-w-[640px] flex-col gap-4 p-4">
          <h1 class="text-3xl font-bold">
            <DynamicInput :modelValue="reaction.name" @update:modelValue="updateName" />
          </h1>

          <Label for="description">Description</Label>
          <GrowTextarea
            id="description"
            :modelValue="reaction.description"
            @update:modelValue="updateDescription"
            placeholder="Enter a description"
          />

          <Card>
            <CardHeader>
              <CardTitle> Trigger </CardTitle>
            </CardHeader>
            <CardContent>
              <Label for="trigger">Type</Label>
              <Select
                id="trigger"
                :modelValue="reaction.trigger"
                @update:modelValue="updateTrigger"
              >
                <SelectTrigger class="font-medium">
                  <SelectValue placeholder="Select a trigger" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="playbook">Playbook</SelectItem>
                    <SelectItem value="ticket">Ticket</SelectItem>
                    <SelectItem value="webhook">Incoming Webhook</SelectItem>
                    <SelectItem value="hook">Internal Hook</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>

              <p v-if="reaction.trigger === 'playbook'" class="py-4 text-sm text-muted-foreground">
                You can add this reaction to a playbooks.
              </p>
              <TriggerTicketType
                v-else-if="reaction.trigger === 'ticket'"
                :modelValue="reaction.triggerdata"
                @update:modelValue="updateTriggerData"
              />
              <TriggerWebhook
                v-else-if="reaction.trigger === 'webhook'"
                :modelValue="reaction.triggerdata"
                @update:modelValue="updateTriggerData"
              />
              <TriggerHook
                v-else-if="reaction.trigger === 'hook'"
                :modelValue="reaction.triggerdata"
                @update:modelValue="updateTriggerData"
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle> Reaction </CardTitle>
            </CardHeader>
            <CardContent>
              <Label>Type</Label>
              <Select :modelValue="reaction.reaction" @update:modelValue="updateReaction">
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

              <ReactionPython
                v-if="reaction.reaction === 'python'"
                :modelValue="reaction.reactiondata"
                @update:modelValue="updateReactionData"
              />
              <ReactionWebhook
                v-else-if="reaction.reaction === 'webhook'"
                :modelValue="reaction.reactiondata"
                @update:modelValue="updateReactionData"
              />
            </CardContent>
          </Card>

          {{ reaction }}
        </div>
      </ScrollArea>
    </div>
  </TanView>
</template>
