<script setup lang="ts">
import { marked } from 'marked'
import { computed } from 'vue'

const props = defineProps<{
  nosyntax?: boolean
  markdown: string
}>()

const html = computed(() => marked(props.markdown))
</script>

<template>
  <div :class="{ prose: !nosyntax, 'dark:prose-invert': !nosyntax }" v-html="html" />
</template>
