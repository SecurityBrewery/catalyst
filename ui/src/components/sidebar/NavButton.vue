<script lang="ts" setup>
import Icon from '@/components/Icon.vue'
import { Button, buttonVariants } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import { cn } from '@/lib/utils'

defineProps<{
  isCollapsed: boolean
  index: number
  title: string
  label?: string
  icon: string
  variant: 'default' | 'ghost'
}>()
</script>

<template>
  <Tooltip v-if="isCollapsed" :key="`1-${index}`" :delay-duration="0">
    <TooltipTrigger as-child>
      <Button
        :class="
          cn(
            buttonVariants({ variant: variant, size: 'icon' }),
            'h-9 w-9',
            variant === 'default' &&
              'dark:bg-muted dark:text-muted-foreground dark:hover:bg-muted dark:hover:text-white'
          )
        "
      >
        <Icon :name="icon" class="size-4" />
        <span class="sr-only">{{ title }}</span>
      </Button>
    </TooltipTrigger>
    <TooltipContent side="right" class="flex items-center gap-4">
      {{ title }}
      <span v-if="label" class="ml-auto text-muted-foreground">
        {{ label }}
      </span>
    </TooltipContent>
  </Tooltip>
  <Button
    v-else
    :key="`2-${index}`"
    :class="
      cn(
        buttonVariants({ variant: variant, size: 'sm' }),
        variant === 'default' &&
          'dark:bg-muted dark:text-white dark:hover:bg-muted dark:hover:text-white',
        'justify-start'
      )
    "
  >
    <Icon :name="icon" class="mr-2 size-4" />
    {{ title }}
    <span
      v-if="label"
      :class="cn('ml-auto', variant === 'default' && 'text-background dark:text-white')"
    >
      {{ label }}
    </span>
  </Button>
</template>
