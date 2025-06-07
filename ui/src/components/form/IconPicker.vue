<script setup lang="ts">
import Icon from '@/components/Icon.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

import { computed, ref } from 'vue'

const modelValue = defineModel<string>({
  default: ''
})
const isOpen = ref(false)
const searchTerm = ref('')

const iconMap = [
  // Priority/Importance
  { name: 'Flame' },
  { name: 'Star' },
  { name: 'Flag' },
  { name: 'Bell' },
  { name: 'BellRing' },
  { name: 'Shield' },
  { name: 'ShieldAlert' },
  { name: 'Siren' },

  // Status/State
  { name: 'HelpCircle' },
  { name: 'Info' },
  { name: 'AlertCircle' },
  { name: 'AlertTriangle' },
  { name: 'CheckCircle' },
  { name: 'XCircle' },
  { name: 'Clock' },
  { name: 'Timer' },

  // Action/Process
  { name: 'Calendar' },
  { name: 'FileText' },
  { name: 'FileQuestion' },
  { name: 'FileWarning' },
  { name: 'User' },

  // Communication
  { name: 'MessageSquare' },
  { name: 'MessageCircle' },
  { name: 'Mail' },

  // System/Technical
  { name: 'Settings' },
  { name: 'Wrench' },
  { name: 'Hammer' },
  { name: 'Code' },
  { name: 'Bug' },
  { name: 'Database' },
  { name: 'Server' },
  { name: 'Cloud' }
]

// Get all icon names that are actual components
const iconNames = computed(() =>
  iconMap.filter((icon) => {
    return icon.name.toLowerCase().includes(searchTerm.value.toLowerCase())
  })
)

const selectIcon = (name: string) => {
  modelValue.value = name
  isOpen.value = false
}
</script>

<template>
  <div class="flex items-center gap-2">
    <Popover>
      <PopoverTrigger>
        <Button
          type="button"
          variant="outline"
          role="combobox"
          class="flex items-center justify-start gap-2"
        >
          <Icon :name="modelValue || 'CircleHelp'" class="size-4" />
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" class="grid w-full grid-cols-8 gap-2 p-2">
        <Button
          v-for="icon in iconNames"
          :key="icon.name"
          type="button"
          @click="selectIcon(icon.name)"
          class="cursor-pointer"
          variant="ghost"
        >
          <Icon :name="icon.name || 'CircleHelp'" class="size-4" />
        </Button>
      </PopoverContent>
    </Popover>
    <Input v-model="modelValue" />
  </div>
</template>
