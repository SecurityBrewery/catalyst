<script lang="ts" setup>
import TanView from '@/components/TanView.vue'
import JSONSchemaFormFields from '@/components/form/JSONSchemaFormFields.vue'
import DynamicMDEditor from '@/components/input/DynamicMDEditor.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import StatusIcon from '@/components/ticket/StatusIcon.vue'
import TicketActionBar from '@/components/ticket/TicketActionBar.vue'
import TicketCloseBar from '@/components/ticket/TicketCloseBar.vue'
import TicketHeader from '@/components/ticket/TicketHeader.vue'
import TicketTab from '@/components/ticket/TicketTab.vue'
import TicketComments from '@/components/ticket/comment/TicketComments.vue'
import TicketFiles from '@/components/ticket/file/TicketFiles.vue'
import TicketLinks from '@/components/ticket/link/TicketLinks.vue'
import TicketTasks from '@/components/ticket/task/TicketTasks.vue'
import TicketTimeline from '@/components/ticket/timeline/TicketTimeline.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { Edit } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket, Type } from '@/lib/types'
import { handleError } from '@/lib/utils'

const route = useRoute()
const queryClient = useQueryClient()

defineProps<{
  selectedType: Type
}>()

const id = computed(() => route.params.id as string)

const message = ref('')
const editMode = ref(false)

const {
  isPending,
  isError,
  data: ticket,
  error
} = useQuery({
  queryKey: ['tickets', id.value],
  queryFn: (): Promise<Ticket> =>
    pb.collection('tickets').getOne(id.value, {
      expand:
        'type,owner,comments_via_ticket.author,files_via_ticket,timeline_via_ticket,links_via_ticket,tasks_via_ticket.owner'
    })
})

const editDescriptionMutation = useMutation({
  mutationFn: () =>
    pb.collection('tickets').update(id.value, {
      description: message.value
    }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', id.value] })
    editMode.value = false
  },
  onError: handleError
})

const edit = () => (editMode.value = true)

const editStateMutation = useMutation({
  mutationFn: (state: Record<string, any>): Promise<Ticket> =>
    pb.collection('tickets').update(id.value, {
      state: state
    }),
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tickets', id.value] }),
  onError: handleError
})

const taskStatus = computed(() => {
  if (!ticket.value) {
    return 'pending'
  }

  const tasks = ticket.value.expand.tasks_via_ticket

  if (tasks.every((task) => !task.open)) {
    return 'completed'
  }

  if (tasks.every((task) => task.open)) {
    return 'open'
  }

  return 'pending'
})

const updateDescription = (value: string) => (message.value = value)
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <template v-if="ticket">
      <TicketActionBar :ticket="ticket" class="shrink-0" />
      <ColumnBody>
        <ColumnBodyContainer class="flex-col gap-4 xl:flex-row">
          <div class="flex flex-1 flex-col gap-4">
            <TicketHeader :ticket="ticket" />
            <Card class="relative p-4">
              <Button
                v-if="!editMode"
                variant="outline"
                class="float-right h-8 gap-2"
                @click="edit"
              >
                <Edit class="h-3.5 w-3.5" />
                <span>Edit</span>
              </Button>
              <DynamicMDEditor
                :modelValue="ticket.description"
                @update:modelValue="updateDescription"
                v-model:edit="editMode"
                autofocus
                placeholder="Type a description..."
                @save="editDescriptionMutation.mutate"
                class="min-h-14"
              />
            </Card>
            <Separator />
            <Tabs default-value="timeline" class="flex flex-1 flex-col">
              <TabsList>
                <TabsTrigger value="timeline">
                  Timeline
                  <Badge
                    v-if="
                      ticket.expand.timeline_via_ticket &&
                      ticket.expand.timeline_via_ticket.length > 0
                    "
                    variant="outline"
                    class="ml-2 hidden sm:inline-flex"
                  >
                    {{
                      ticket.expand.timeline_via_ticket
                        ? ticket.expand.timeline_via_ticket.length
                        : 0
                    }}
                  </Badge>
                </TabsTrigger>
                <TabsTrigger value="tasks">
                  Tasks
                  <Badge
                    v-if="
                      ticket.expand.tasks_via_ticket && ticket.expand.tasks_via_ticket.length > 0
                    "
                    variant="outline"
                    class="ml-2 hidden sm:inline-flex"
                  >
                    {{ ticket.expand.tasks_via_ticket ? ticket.expand.tasks_via_ticket.length : 0 }}
                    <StatusIcon :status="taskStatus" class="size-6" />
                  </Badge>
                </TabsTrigger>
                <TabsTrigger value="comments">
                  Comments
                  <Badge
                    v-if="
                      ticket.expand.comments_via_ticket &&
                      ticket.expand.comments_via_ticket.length > 0
                    "
                    variant="outline"
                    class="ml-2 hidden sm:inline-flex"
                  >
                    {{
                      ticket.expand.comments_via_ticket
                        ? ticket.expand.comments_via_ticket.length
                        : 0
                    }}
                  </Badge>
                </TabsTrigger>
              </TabsList>
              <TicketTab value="timeline">
                <TicketTimeline :ticket="ticket" :timeline="ticket.expand.timeline_via_ticket" />
              </TicketTab>
              <TicketTab value="tasks">
                <TicketTasks :ticket="ticket" :tasks="ticket.expand.tasks_via_ticket" />
              </TicketTab>
              <TicketTab value="comments">
                <TicketComments :ticket="ticket" :comments="ticket.expand.comments_via_ticket" />
              </TicketTab>
            </Tabs>
            <Separator class="xl:hidden" />
          </div>
          <div class="flex flex-col gap-4 xl:w-96 xl:flex-initial">
            <div>
              <div class="flex h-10 flex-row items-center justify-between text-muted-foreground">
                <span class="text-sm font-semibold"> Details </span>
              </div>
              <JSONSchemaFormFields
                :modelValue="ticket.state"
                @update:modelValue="editStateMutation.mutate"
                :schema="selectedType.schema"
              />
            </div>
            <Separator />
            <TicketLinks :ticket="ticket" :links="ticket.expand.links_via_ticket" />
            <Separator />
            <TicketFiles :ticket="ticket" :files="ticket.expand.files_via_ticket" />
          </div>
        </ColumnBodyContainer>
      </ColumnBody>
      <Separator />
      <TicketCloseBar :ticket="ticket" class="shrink-0" />
    </template>
  </TanView>
</template>
