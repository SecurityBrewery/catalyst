<script setup lang="ts">
import { Button } from '@/components/ui/button'

import { ChevronsRight } from 'lucide-vue-next'

import { PaginationLast, type PaginationLastProps } from 'radix-vue'
import { type HTMLAttributes, computed } from 'vue'

import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<PaginationLastProps & { class?: HTMLAttributes['class'] }>(),
  {
    asChild: true
  }
)

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
  <PaginationLast v-bind="delegatedProps">
    <Button :class="cn('h-10 w-10 p-0', props.class)" variant="outline">
      <slot>
        <ChevronsRight class="h-4 w-4" />
      </slot>
    </Button>
  </PaginationLast>
</template>
