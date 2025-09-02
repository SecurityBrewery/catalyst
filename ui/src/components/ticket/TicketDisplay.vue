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
import { useToast } from '@/components/ui/toast/use-toast'

import { Edit } from 'lucide-vue-next'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

import { useAPI } from '@/api'
import type {
  ExtendedComment,
  ExtendedTask,
  ExtendedTicket,
  Link,
  ModelFile,
  Ticket,
  TimelineEntry,
  Type
} from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()

const route = useRoute()
const queryClient = useQueryClient()
const { toast } = useToast()

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
  queryFn: (): Promise<ExtendedTicket> => api.getTicket({ id: id.value })
})

const { data: timeline } = useQuery({
  queryKey: ['timeline', id.value],
  queryFn: (): Promise<Array<TimelineEntry>> => api.listTimeline({ ticket: id.value })
})

const { data: tasks } = useQuery({
  queryKey: ['tasks', id.value],
  queryFn: (): Promise<Array<ExtendedTask>> => api.listTasks({ ticket: id.value })
})

const { data: comments } = useQuery({
  queryKey: ['comments', id.value],
  queryFn: (): Promise<Array<ExtendedComment>> => api.listComments({ ticket: id.value })
})

const { data: files } = useQuery({
  queryKey: ['files', id.value],
  queryFn: (): Promise<Array<ModelFile>> => api.listFiles({ ticket: id.value })
})

const { data: links } = useQuery({
  queryKey: ['links', id.value],
  queryFn: (): Promise<Array<Link>> => api.listLinks({ ticket: id.value })
})

const editDescriptionMutation = useMutation({
  mutationFn: () =>
    api.updateTicket({ id: id.value, ticketUpdate: { description: message.value } }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', id.value] })
    toast({
      title: 'Ticket updated',
      description: 'The ticket description has been updated'
    })
    editMode.value = false
  },
  onError: handleError('Failed to update description')
})

const edit = () => (editMode.value = true)

const editStateMutation = useMutation({
  mutationFn: (state: Record<string, any>): Promise<Ticket> =>
    api.updateTicket({ id: id.value, ticketUpdate: { state } }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['tickets', id.value] })
    toast({
      title: 'Ticket updated',
      description: 'The ticket state has been updated'
    })
  },
  onError: handleError('Failed to update state')
})

const taskStatus = computed(() => {
  if (!ticket.value) {
    return 'pending'
  }

  if (tasks.value && tasks.value.every((task) => !task.open)) {
    return 'completed'
  }

  if (tasks.value && tasks.value.every((task) => task.open)) {
    return 'open'
  }

  return 'pending'
})

const updateDescription = (value: string | undefined) => (message.value = value ?? '')
</script>

<template>
  <TanView :isError="isError" :isPending="isPending" :error="error">
    <template v-if="ticket">
      <TicketActionBar :ticket="ticket" />
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
                    v-if="timeline && timeline.length > 0"
                    variant="outline"
                    class="ml-2 hidden sm:inline-flex"
                  >
                    {{ timeline.length }}
                  </Badge>
                </TabsTrigger>
                <TabsTrigger value="tasks">
                  Tasks
                  <Badge
                    v-if="tasks && tasks.length > 0"
                    variant="outline"
                    class="ml-2 hidden sm:inline-flex"
                  >
                    {{ tasks.length }}
                    <StatusIcon :status="taskStatus" class="size-6" />
                  </Badge>
                </TabsTrigger>
                <TabsTrigger value="comments">
                  Comments
                  <Badge
                    v-if="comments && comments.length > 0"
                    variant="outline"
                    class="ml-2 hidden sm:inline-flex"
                  >
                    {{ comments.length }}
                  </Badge>
                </TabsTrigger>
              </TabsList>
              <TicketTab value="timeline">
                <TicketTimeline :ticket="ticket" :timeline="timeline" />
              </TicketTab>
              <TicketTab value="tasks">
                <TicketTasks :ticket="ticket" :tasks="tasks" />
              </TicketTab>
              <TicketTab value="comments">
                <TicketComments :ticket="ticket" :comments="comments" />
              </TicketTab>
            </Tabs>
            <Separator class="xl:hidden" />
          </div>
          <div class="flex flex-col gap-4 xl:w-96 xl:flex-initial">
            <div>
              <div class="flex h-10 flex-row items-center justify-between">
                <span class="text-sm font-medium"> Details </span>
              </div>
              <JSONSchemaFormFields
                :modelValue="ticket.state"
                @update:modelValue="editStateMutation.mutate"
                :schema="selectedType.schema"
              />
            </div>
            <Separator />
            <TicketLinks :ticket="ticket" :links="links" />
            <Separator />
            <TicketFiles :ticket="ticket" :files="files" />
          </div>
        </ColumnBodyContainer>
      </ColumnBody>
      <Separator />
      <TicketCloseBar :ticket="ticket" />
    </template>
  </TanView>
</template>
