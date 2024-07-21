<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import ReactionDisplay from '@/components/reaction/ReactionDisplay.vue'
import ReactionList from '@/components/reaction/ReactionList.vue'
import ReactionNew from '@/components/reaction/ReactionNew.vue'

import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'

const route = useRoute()
const router = useRouter()

const id = computed(() => route.params.id as string)

onMounted(() => {
  if (!pb.authStore.model) {
    router.push({ name: 'login' })
  }
})
</script>

<template>
  <ThreeColumn>
    <template #list>
      <ReactionList />
    </template>
    <template #single>
      <div v-if="!id" class="flex h-full w-full items-center justify-center text-lg text-gray-500">
        No reaction selected
      </div>
      <ReactionNew v-else-if="id === 'new'" key="new" />
      <ReactionDisplay v-else :key="id" :id="id" />
    </template>
  </ThreeColumn>
</template>
