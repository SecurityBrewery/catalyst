<script setup lang="ts">
import { CommandEmpty, CommandGroup, CommandItem, CommandList } from '@/components/ui/command'
import {
  TagsInput,
  TagsInputInput,
  TagsInputItem,
  TagsInputItemDelete,
  TagsInputItemText
} from '@/components/ui/tags-input'

import { ComboboxAnchor, ComboboxInput, ComboboxPortal, ComboboxRoot } from 'radix-vue'
import { computed, ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string[]
    items: string[]
    placeholder?: string
  }>(),
  {
    modelValue: () => [],
    items: () => [],
    placeholder: ''
  }
)

const emit = defineEmits(['update:modelValue'])

const open = ref(false)
const searchTerm = ref('')
const selectedItems = ref<string[]>(props.modelValue)

watch(
  () => selectedItems.value,
  (value) => emit('update:modelValue', value),
  { deep: true }
)

const filteredItems = computed(() => {
  if (!selectedItems.value) return props.items
  return props.items.filter((i) => !selectedItems.value.includes(i))
})
</script>

<template>
  <TagsInput class="flex items-center gap-2 px-0" :modelValue="selectedItems">
    <div class="flex flex-wrap items-center">
      <TagsInputItem v-for="item in selectedItems" :key="item" :value="item" class="ml-2">
        <TagsInputItemText />
        <TagsInputItemDelete />
      </TagsInputItem>
    </div>

    <ComboboxRoot
      v-model="selectedItems"
      v-model:open="open"
      v-model:searchTerm="searchTerm"
      class="flex-1"
    >
      <ComboboxAnchor as-child>
        <ComboboxInput
          :placeholder="placeholder"
          as-child
          :class="selectedItems.length < items.length ? '' : 'hidden'"
        >
          <TagsInputInput @keydown.enter.prevent @focus="open = true" @blur="open = false" />
        </ComboboxInput>
      </ComboboxAnchor>

      <ComboboxPortal>
        <CommandList
          v-if="selectedItems.length < items.length"
          position="popper"
          class="mt-2 w-[--radix-popper-anchor-width] rounded-md border bg-popover text-popover-foreground shadow-md outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"
        >
          <CommandEmpty />
          <CommandGroup>
            <CommandItem
              v-for="item in filteredItems"
              :key="item"
              :value="item"
              @select.prevent="
                (ev) => {
                  if (typeof ev.detail.value === 'string') {
                    searchTerm = ''
                    selectedItems.push(ev.detail.value)
                  }

                  if (filteredItems.length === 0) {
                    open = false
                  }
                }
              "
            >
              {{ item }}
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </ComboboxPortal>
    </ComboboxRoot>
  </TagsInput>
</template>
