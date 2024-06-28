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

import { LoaderCircle, Search } from 'lucide-vue-next'

import { useQuery } from '@tanstack/vue-query'
import debounce from 'lodash.debounce'
import type { ListResult } from 'pocketbase'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'
import type { Ticket, Type } from '@/lib/types'

const router = useRouter()

const props = defineProps<{
  selectedType: Type
}>()

const searchValue = ref('')
const tab = ref('open')

const filter = computed(() => {
  let raw = ''
  const params: Record<string, string> = {}

  /*
  if (searchValue.value && searchValue.value !== '') {
    let raws: Array<string> = []
    props.selectedType.expand.fields.forEach((field) => {
      if (field.type === 'bool') return
      raws.push(`${field.name} ~ {:search}`)
    })
    raw += '(' + raws.join(' || ') + `)`
    params['search'] = searchValue.value
  }
   */

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
      expand: 'owner,type'
    })
})

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
          <div class="relative">
            <Search class="absolute left-2 top-2.5 size-4 text-muted-foreground" />
            <Input v-model="searchValue" placeholder="Search" class="pl-8" />
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
      <!-- TabsContent value="all" class="flex flex-1 flex-col items-center ">
          <TicketListList :tickets="tickets" />
        </TabsContent>
        <TabsContent value="open" class="flex flex-1 flex-col items-center ">
          <TicketListList :tickets="tickets" />
        </TabsContent>
        <TabsContent value="closed" class="flex flex-1 flex-col items-center">
          <TicketListList :tickets="tickets" />
        </TabsContent-->
    </Tabs>
  </div>
</template>

<style scoped>
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(15px);
}

.list-leave-active {
  position: absolute;
}
</style>
