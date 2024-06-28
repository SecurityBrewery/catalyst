<script lang="ts" setup>
import Icon from '@/components/Icon.vue'
import { buttonVariants } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import { computed } from 'vue'
import { useRoute } from 'vue-router'

import { cn } from '@/lib/utils'

const route = useRoute()

export interface LinkProp {
  title: string
  label?: string
  icon: string
  to: string
  variant: 'default' | 'ghost'
  disabled?: boolean
}

const props = defineProps<{
  isCollapsed: boolean
  index: number
  link: LinkProp
}>()

const variant = computed((): 'default' | 'ghost' =>
  route.path.startsWith(props.link.to) ? 'default' : 'ghost'
)
</script>

<template>
  <Tooltip v-if="isCollapsed" :key="`1-${index}`" :delay-duration="0">
    <TooltipTrigger as-child>
      <component
        :is="link.disabled ? 'span' : 'router-link'"
        :to="link.to"
        :class="
          cn(
            buttonVariants({ variant: variant, size: 'icon' }),
            'h-9 w-9',
            link.variant === 'default' &&
              'dark:bg-muted dark:text-muted-foreground dark:hover:bg-muted dark:hover:text-white',
            link.disabled &&
              'text-muted-foreground hover:bg-transparent hover:text-muted-foreground'
          )
        "
      >
        <Icon :name="link.icon" class="size-4" />
        <span class="sr-only">{{ link.title }}</span>
      </component>
    </TooltipTrigger>
    <TooltipContent side="right" class="flex items-center gap-4">
      {{ link.title }}
      <span v-if="link.label" class="ml-auto text-muted-foreground">
        {{ link.label }}
      </span>
    </TooltipContent>
  </Tooltip>

  <component
    :is="link.disabled ? 'span' : 'router-link'"
    v-else
    :key="`2-${index}`"
    :to="link.to"
    :class="
      cn(
        buttonVariants({ variant: variant, size: 'sm' }),
        link.variant === 'default' &&
          'dark:bg-muted dark:text-white dark:hover:bg-muted dark:hover:text-white',
        'justify-start',
        link.disabled && 'text-muted-foreground hover:bg-transparent hover:text-muted-foreground'
      )
    "
  >
    <Icon :name="link.icon" class="mr-2 size-4" />
    {{ link.title }}
    <span
      v-if="link.label"
      :class="cn('ml-auto', link.variant === 'default' && 'text-background dark:text-white')"
    >
      {{ link.label }}
    </span>
  </component>
</template>
