<script lang="ts" setup>
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
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
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { LoaderCircle, Search } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import debounce from 'lodash.debounce'
import type { ListResult } from 'pocketbase'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { SearchTicket, Type } from '@/lib/types'

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
      'owner_name ~ {:search}',
      'comment_messages ~ {:search}',
      'file_names ~ {:search}',
      'link_names ~ {:search}',
      'link_urls ~ {:search}',
      'task_names ~ {:search}',
      'timeline_messages ~ {:search}'
    ]

    Object.keys(props.selectedType.schema.properties).forEach((key) => {
      const property = props.selectedType.schema.properties[key]
      if (property.type === 'bool') return
      raws.push(`state.${key} ~ {:search}`)
    })
    raw += '(' + raws.join(' || ') + `)`
    params['search'] = searchValue.value
  }

  if (tab.value === 'open') {
    if (raw !== '') raw += ' && '
    raw += 'open = true'
  } else if (tab.value === 'closed') {
    if (raw !== '') raw += ' && '
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
  queryFn: (): Promise<ListResult<SearchTicket>> =>
    pb.collection('ticket_search').getList(page.value, perPage.value, {
      sort: '-created',
      filter: filter.value
    })
})

watch(
  () => ticketItems.value,
  () => {
    if (!route.params.id && ticketItems.value && ticketItems.value.items.length > 0) {
      router.push({
        name: 'tickets',
        params: { type: props.selectedType.id }
      })
    }
  }
)

const debouncedRefetch = debounce(refetch, 300)
watch(
  () => searchValue.value,
  () => debouncedRefetch()
)
watch([tab, props.selectedType, page, perPage], () => refetch())
</script>

<template>
  <ColumnHeader :title="selectedType?.plural">
    <div class="ml-auto">
      <TicketNewDialog :selectedType="selectedType" />
    </div>
  </ColumnHeader>
  <Tabs v-model="tab" class="flex flex-1 flex-col overflow-hidden">
    <div class="flex items-center justify-between px-2 pt-2">
      <TabsList>
        <TabsTrigger value="all">All</TabsTrigger>
        <TabsTrigger value="open">Open</TabsTrigger>
        <TabsTrigger value="closed">Closed</TabsTrigger>
      </TabsList>
      <!-- Button variant="outline" size="sm" class="h-7 gap-1 rounded-md px-3">
        <ListFilter class="h-3.5 w-3.5" />
        <span class="sr-only sm:not-sr-only">Filter</span>
      </Button-->
    </div>
    <div class="p-2">
      <form>
        <div class="relative flex flex-row items-center">
          <Input v-model="searchValue" placeholder="Search" @keydown.enter.prevent class="pl-8" />
          <span class="absolute inset-y-0 start-0 flex items-center justify-center px-2">
            <Search class="size-4 text-muted-foreground" />
          </span>
        </div>
      </form>
    </div>
    <Separator />
    <div v-if="isPending" class="flex h-full w-full items-center justify-center">
      <LoaderCircle class="h-16 w-16 animate-spin text-primary" />
    </div>
    <Alert v-else-if="isError" variant="destructive" class="mb-2 h-screen w-screen">
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <div v-else-if="ticketItems" class="flex-1 overflow-y-auto overflow-x-hidden">
      <TicketListList :tickets="ticketItems.items" />
    </div>
    <Separator />
    <div class="my-2 flex items-center justify-center">
      <span class="text-xs text-muted-foreground">
        {{ ticketItems ? ticketItems.items.length : '?' }} of
        {{ ticketItems ? ticketItems.totalItems : '?' }} tickets
      </span>
    </div>
    <div class="mb-2 flex items-center justify-center">
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
              <Button class="h-10 w-10 p-0" :variant="item.value === page ? 'default' : 'outline'">
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
</template>
