<script setup lang="ts">
import Icon from '@/components/Icon.vue'
import type { Resource, Service } from '@/components/editor/common/enrichtments'

import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

import { pb } from '@/lib/pocketbase'

const props = defineProps<{
  resource: Resource
}>()

const { data: services } = useQuery({
  queryKey: ['services'],
  queryFn: (): Promise<Array<Service>> =>
    pb
      .send('/api/analysis/services', {})
      .then((data) => {
        return data.services
      })
      .catch((err) => {
        return Promise.reject(err)
      })
})

const resourceTypeName = computed(() => {
  if (services.value) {
    for (const service of services.value) {
      if (service.id === props.resource.service) {
        for (const resource_type of service.resource_types) {
          if (resource_type.id === props.resource.type) {
            return resource_type.name
          }
        }
      }
    }
  }

  return props.resource.type
})
</script>

<template>
  <div class="flex flex-col">
    <div class="flex flex-row items-center gap-2">
      <Icon v-if="resource.icon" :name="resource.icon" class="size-3.5 flex-shrink-0" />
      <span class="text-xs text-muted-foreground">
        {{ resourceTypeName }}
      </span>
    </div>
    <h2 class="mt-0.5 text-sm font-semibold">
      {{ resource.name }}
    </h2>
    <div
      v-if="resource.attributes && resource.attributes.length > 0"
      class="mt-1 flex flex-col gap-1 text-xs"
    >
      <div v-for="attribute in resource.attributes" :key="attribute.title">
        <div class="flex flex-row items-center gap-1">
          <Icon :name="attribute.icon" :stroke-width="2" class="size-3 text-muted-foreground" />
          <strong class="text-muted-foreground">
            {{ attribute.title }}
          </strong>
          <span>
            {{ attribute.value }}
          </span>
        </div>
      </div>
    </div>
    <div class="mt-1 text-sm text-muted-foreground">
      {{
        resource.description && resource.description.length > 160
          ? resource.description.slice(0, 150) + '...'
          : resource.description
      }}
    </div>
    <div class="mt-1 flex flex-row gap-2">
      <a
        :href="resource.url"
        target="_blank"
        rel="noopener noreferrer"
        class="text-xs text-muted-foreground"
      >
        {{ resource.url }}
      </a>
    </div>
  </div>
</template>
