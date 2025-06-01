<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import UserDisplay from '@/components/user/UserDisplay.vue'
import UserList from '@/components/user/UserList.vue'
import UserNew from '@/components/user/UserNew.vue'

import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const id = computed(() => route.params.id as string)
</script>

<template>
  <ThreeColumn :show-details="!!id">
    <template #list>
      <UserList />
    </template>
    <template #single>
      <ColumnBody v-if="!id" class="items-center justify-center text-lg text-gray-500">
        No user selected
      </ColumnBody>
      <UserNew v-else-if="id === 'new'" key="new" />
      <UserDisplay v-else :key="id" :id="id" />
    </template>
  </ThreeColumn>
</template>
