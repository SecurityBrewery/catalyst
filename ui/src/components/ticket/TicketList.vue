<script lang="ts" setup>
import TicketListList from '@/components/ticket/TicketListList.vue'
import TicketNewDialog from '@/components/ticket/TicketNewDialog.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Pagination,
  PaginationEllipsis,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationListItem,
  PaginationNext,
  PaginationPrev
} from '@/components/ui/pagination'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { Info, LoaderCircle, Search } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import debounce from 'lodash.debounce'
import type { ListResult } from 'pocketbase'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket, Type } from '@/lib/types'

const router = useRouter()
const route = useRoute()

const props = defineProps<{
  selectedType: Type
}>()

const searchValue = ref('')
const tab = ref('open')

const filter = computed(() => {
  let raw = ''
  const params: Record<string, string> = {}

  if (searchValue.value && searchValue.value !== '') {
    let raws: Array<string> = [
      'name ~ {:search}',
      'description ~ {:search}',
      'owner.name ~ {:search}',
      'owner.email ~ {:search}',
      'links_via_ticket.name ~ {:search}',
      'links_via_ticket.url ~ {:search}',
      'tasks_via_ticket.name ~ {:search}',
      'comments_via_ticket.message ~ {:search}',
      'files_via_ticket.name ~ {:search}',
      'timeline_via_ticket.message ~ {:search}',
      'state.severity ~ {:search}'
    ]

    Object.keys(props.selectedType.schema.properties).forEach((key) => {
      const property = props.selectedType.schema.properties[key]
      if (property.type === 'bool') return
      raws.push(`state.${key} ~ {:search}`)
    })
    raw += '(' + raws.join(' || ') + `)`
    params['search'] = searchValue.value
  }

  if (raw !== '') raw += ' && '
  if (tab.value === 'open') {
    raw += 'open = true'
  } else if (tab.value === 'closed') {
    raw += 'open = false'
  }

  if (raw !== '') raw += ' && '
  raw += 'type = {:type}'
  params['type'] = props.selectedType.id

  if (raw === '') return ''

  return pb.filter(raw, params)
})

const page = ref(1)
const perPage = ref(10)

const {
  isPending,
  isError,
  data: ticketItems,
  error,
  refetch
} = useQuery({
  queryKey: ['tickets', filter.value],
  queryFn: (): Promise<ListResult<Ticket>> =>
    pb.collection('tickets').getList(page.value, perPage.value, {
      sort: '-created',
      filter: filter.value,
      expand:
        'type,owner,comments_via_ticket.author,files_via_ticket,timeline_via_ticket,links_via_ticket,tasks_via_ticket.owner'
    })
})

watch(
  () => ticketItems.value,
  () => {
    if (!route.params.id && ticketItems.value && ticketItems.value.items.length > 0) {
      router.push({
        name: 'tickets',
        params: { type: props.selectedType.id, id: ticketItems.value.items[0].id }
      })
    }
  }
)

const debouncedRefetch = debounce(refetch, 300)
watch(searchValue, () => debouncedRefetch())
watch([tab, props.selectedType, page, perPage], () => refetch())
</script>

<template>
  <div class="flex h-screen flex-col">
    <div class="flex items-center bg-background px-4 py-2">
      <h1 class="text-xl font-bold">
        {{ selectedType?.plural }}
      </h1>
      <div class="ml-auto">
        <TicketNewDialog :selectedType="selectedType" />
      </div>
    </div>
    <Separator />
    <Tabs v-model="tab" class="flex flex-1 flex-col overflow-hidden">
      <div class="flex items-center justify-between px-4 pt-2">
        <TabsList>
          <TabsTrigger value="all"> All</TabsTrigger>
          <TabsTrigger value="open"> Open</TabsTrigger>
          <TabsTrigger value="closed"> Closed</TabsTrigger>
        </TabsList>
        <!-- Button variant="outline" size="sm" class="h-7 gap-1 rounded-md px-3">
          <ListFilter class="h-3.5 w-3.5" />
          <span class="sr-only sm:not-sr-only">Filter</span>
        </Button-->
      </div>
      <div class="px-4 py-2">
        <form>
          <div class="relative flex flex-row items-center">
            <Input v-model="searchValue" placeholder="Search" class="pl-8" />
            <span class="absolute inset-y-0 start-0 flex items-center justify-center px-2">
              <Search class="size-4 text-muted-foreground" />
            </span>

            <div>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Info class="ml-2 size-4 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent>
                    <p class="w-64">
                      Search name, description, or owner. Links, tasks, comments, files, and
                      timeline messages are also searched, but cause unreliable results if there are
                      more than 1000 records.
                    </p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </div>
        </form>
      </div>
      <Separator />
      <div v-if="isPending" class="flex h-full w-full items-center justify-center">
        <LoaderCircle class="h-16 w-16 animate-spin text-primary" />
      </div>
      <Alert v-else-if="isError" variant="destructive" class="mb-4 h-screen w-screen">
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>
      <ScrollArea v-else-if="ticketItems" class="flex-1">
        <TicketListList :tickets="ticketItems.items" />
      </ScrollArea>
      <Separator />
      <div class="my-2 flex items-center justify-center">
        <span class="text-xs text-muted-foreground">
          {{ ticketItems ? ticketItems.items.length : '?' }} of
          {{ ticketItems ? ticketItems.totalItems : '?' }} tickets
        </span>
      </div>
      <div class="mb-4 flex items-center justify-center">
        <Pagination
          v-slot="{ page }"
          :total="ticketItems ? ticketItems.totalItems : 0"
          :itemsPerPage="perPage"
          :sibling-count="0"
          :default-page="1"
          @update:page="page = $event"
        >
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem
                v-if="item.type === 'page'"
                :key="index"
                :value="item.value"
                as-child
              >
                <Button
                  class="h-10 w-10 p-0"
                  :variant="item.value === page ? 'default' : 'outline'"
                >
                  {{ item.value }}
                </Button>
              </PaginationListItem>
              <PaginationEllipsis v-else :key="item.type" :index="index" />
            </template>

            <PaginationNext />
            <PaginationLast />
          </PaginationList>
        </Pagination>
      </div>
    </Tabs>
  </div>
</template>