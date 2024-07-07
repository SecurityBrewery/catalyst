<script setup lang="ts">
import MDEditor from '@/components/input/MDEditor.vue'
import MarkdownView from '@/components/input/MarkdownView.vue'

const model = defineModel<string>()

const edit = defineModel('edit', {
  type: Boolean,
  required: false,
  default: false
})

const emit = defineEmits(['save'])

export interface Props {
  placeholder?: string
  hideCancel?: boolean
  autofocus?: boolean
}

withDefaults(defineProps<Props>(), {
  placeholder: '',
  hideCancel: false,
  autofocus: false
})

const cancel = () => (edit.value = false)
</script>

<template>
  <div class="mt-1">
    <MDEditor
      v-if="edit"
      v-model="model"
      class="-mx-1 -mt-1"
      :autofocus="autofocus"
      :placeholder="placeholder"
      :hideCancel="hideCancel"
      @save="emit('save')"
      @cancel="cancel"
    />

    <MarkdownView v-else :markdown="model ? model : ''" />
  </div>
</template>
