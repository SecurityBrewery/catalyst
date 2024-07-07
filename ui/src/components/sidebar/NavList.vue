<script lang="ts" setup>
import NavItem from '@/components/sidebar/NavItem.vue'

export interface LinkProp {
  title: string
  label?: string
  icon: string
  to: string
  variant: 'default' | 'ghost'
  disabled?: boolean
}

defineProps<{
  isCollapsed: boolean
  links: LinkProp[]
}>()
</script>

<template>
  <div
    :data-collapsed="isCollapsed"
    class="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2"
  >
    <nav
      class="grid gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2"
    >
      <template v-for="(link, index) of links" :key="index">
        <NavItem :index="index" :link="link" :is-collapsed="isCollapsed" />
      </template>
    </nav>
  </div>
</template>
