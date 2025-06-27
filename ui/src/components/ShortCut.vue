<script setup lang="ts">
defineProps<{
  keys: string
}>()
</script>

<template>
  <kbd
    class="pointer-events-none ml-2 inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100"
  >
    <span class="text-xs">
      {{ keys }}
    </span>
  </kbd>
</template>
