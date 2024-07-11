<script setup lang="ts">
import ThreeColumn from '@/components/layout/ThreeColumn.vue'
import ReactionList from '@/components/reaction/ReactionList.vue'
import ReactionPythonDisplay from '@/components/reaction/ReactionPythonDisplay.vue'
import ReactionWebhookDisplay from '@/components/reaction/ReactionWebhookDisplay.vue'

import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { pb } from '@/lib/pocketbase'

const route = useRoute()
const router = useRouter()

const id = computed(() => route.params.id as string)
const type = computed(() => route.params.type as string)

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
      <ReactionPythonDisplay v-else-if="type == 'python'" :key="'python_' + id" :id="id" />
      <ReactionWebhookDisplay v-else-if="type == 'webhook'" :key="'webhook_' + id" :id="id" />
    </template>
  </ThreeColumn>
</template>
