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
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAPI } from '@/api'
import type { TicketSearch } from '@/client/models'
import type { Type } from '@/client/models/Type'

const api = useAPI()

const router = useRouter()
const route = useRoute()

const props = defineProps<{
  selectedType: Type
}>()

const searchValue = ref('')
const tab = ref('open')

const page = ref(1)
const perPage = ref(10)

const totalItems = ref(0)
const offset = computed(() => (page.value - 1) * perPage.value)
const paginationHuman = computed(() => {
  if (totalItems.value === 0) {
    return '0 tickets'
  }

  if (offset.value + 1 === totalItems.value) {
    return `${offset.value + 1} of ${totalItems.value} tickets`
  }

  if (offset.value + perPage.value >= totalItems.value) {
    return `${offset.value + 1} - ${totalItems.value} of ${totalItems.value} tickets`
  }

  return `${offset.value + 1} - ${offset.value + perPage.value} of ${totalItems.value} tickets`
})

const {
  isPending,
  isError,
  data: ticketItems,
  error,
  refetch
} = useQuery({
  queryKey: [
    'tickets',
    searchValue.value,
    tab.value,
    props.selectedType.id,
    page.value,
    perPage.value
  ],
  queryFn: async (): Promise<Array<TicketSearch>> => {
    const response = await api.searchTicketsRaw({
      type: props.selectedType.id,
      query: searchValue.value,
      open: tab.value === 'open' ? true : tab.value === 'closed' ? false : undefined,
      offset: offset.value,
      limit: perPage.value
    })
    totalItems.value = parseInt(response.raw.headers.get('X-Total-Count') ?? '0')

    return response.value()
  }
})

watch(
  () => ticketItems.value,
  () => {
    if (!route.params.id && ticketItems.value && ticketItems.value.length > 0) {
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
            <Search class="text-muted-foreground size-4" />
          </span>
        </div>
      </form>
    </div>
    <Separator />
    <div v-if="isPending" class="flex h-full w-full items-center justify-center">
      <LoaderCircle class="text-primary h-16 w-16 animate-spin" />
    </div>
    <Alert v-else-if="isError" variant="destructive" class="mb-2 h-screen w-screen">
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <div v-else-if="ticketItems" class="flex-1 overflow-x-hidden overflow-y-auto">
      <TicketListList :tickets="ticketItems" />
    </div>
    <Separator />
    <div class="my-2 flex items-center justify-center">
      <span class="text-muted-foreground text-xs">
        {{ paginationHuman }}
      </span>
    </div>
    <div class="mb-2 flex items-center justify-center">
      <Pagination
        v-slot="{ page }"
        :total="totalItems"
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
                :variant="item.value === page ? 'default' : 'outline-solid'"
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
</template>
