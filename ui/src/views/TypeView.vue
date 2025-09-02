<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import TypeDisplay from '@/components/type/TypeDisplay.vue'
import TypeList from '@/components/type/TypeList.vue'
import TypeNew from '@/components/type/TypeNew.vue'

import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const id = computed(() => route.params.id as string)
</script>

<template>
  <ThreeColumn :show-details="!!id">
    <template #list>
      <TypeList />
    </template>
    <template #single>
      <ColumnBody v-if="!id" class="items-center justify-center text-lg text-gray-500">
        No type selected
      </ColumnBody>
      <TypeNew v-else-if="id === 'new'" key="new" />
      <TypeDisplay v-else :key="id" :id="id" />
    </template>
  </ThreeColumn>
</template>
