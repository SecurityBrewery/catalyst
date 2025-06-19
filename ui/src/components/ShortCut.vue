<script setup lang="ts">
defineProps<{
  keys: string
}>()
</script>

<template>
  <kbd
    class="bg-muted text-muted-foreground pointer-events-none ml-2 inline-flex h-5 items-center gap-1 rounded border px-1.5 font-mono text-[10px] font-medium opacity-100 select-none"
  >
    <span class="text-xs">
      {{ keys }}
    </span>
  </kbd>
</template>
