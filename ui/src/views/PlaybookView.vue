<script setup lang="ts">
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import PlaybookDisplay from '@/components/playbook/PlaybookDisplay.vue'
import PlaybookList from '@/components/playbook/PlaybookList.vue'

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
      <PlaybookList />
    </template>
    <template #single>
      <div v-if="!id" class="flex h-full w-full items-center justify-center text-lg text-gray-500">
        No playbook selected
      </div>
      <PlaybookDisplay v-else :key="id" :id="id" />
    </template>
  </ThreeColumn>
</template>
