<script setup lang="ts">
import GroupDisplay from '@/components/group/GroupDisplay.vue'
import GroupList from '@/components/group/GroupList.vue'
import GroupNew from '@/components/group/GroupNew.vue'
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ThreeColumn from '@/components/layout/ThreeColumn.vue'

import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const id = computed(() => route.params.id as string)
</script>

<template>
  <ThreeColumn :show-details="!!id">
    <template #list>
      <GroupList />
    </template>
    <template #single>
      <ColumnBody v-if="!id" class="items-center justify-center text-lg text-gray-500">
        No user selected
      </ColumnBody>
      <GroupNew v-else-if="id === 'new'" key="new" />
      <GroupDisplay v-else :key="id" :id="id" />
    </template>
  </ThreeColumn>
</template>
