<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import RoleDisplay from '@/components/role/RoleDisplay.vue'
import RoleList from '@/components/role/RoleList.vue'
import RoleNew from '@/components/role/RoleNew.vue'

import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const id = computed(() => route.params.id as string)
</script>

<template>
  <ThreeColumn :show-details="!!id">
    <template #list>
      <RoleList />
    </template>
    <template #single>
      <ColumnBody v-if="!id" class="items-center justify-center text-lg text-gray-500">
        No user selected
      </ColumnBody>
      <RoleNew v-else-if="id === 'new'" key="new" />
      <RoleDisplay v-else :key="id" :id="id" />
    </template>
  </ThreeColumn>
</template>
