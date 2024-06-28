<script setup lang="ts">
import ShortCut from '@/components/ShortCut.vue'
import { Button } from '@/components/ui/button'

import { ref } from 'vue'
import { EditorInstance } from 'vue3-easymde'

const model = defineModel<string>()

const emit = defineEmits(['save', 'cancel'])

export interface Props {
  placeholder?: string
  hideCancel?: boolean
  autofocus?: boolean
}

const focus = ref(false)

const props = withDefaults(defineProps<Props>(), {
  placeholder: '',
  hideCancel: false,
  autofocus: false
})

// meanwhile, vue3-easymde also expose particular instance
const editorInstance = ref<EditorInstance | null>(null)
// you can call getMDEInstance method to get easymde instance
// if (editorInstance.value) {
//   console.log(editorInstance.value.getMDEInstance())
// }
</script>

<template>
  <div>
    <!-- https://github.com/Ionaru/easy-markdown-editor -->
    <vue-easymde
      v-model="model"
      class="prose dark:prose-invert"
      ref="editorInstance"
      @keydown.meta.enter="emit('save')"
      @keydown.ctrl.enter="emit('save')"
      @click="focus = true"
      @blur="focus = false"
      :options="{
        autofocus: props.autofocus,
        placeholder: props.placeholder,
        minHeight: '40px',
        toolbar: false,
        status: false,
        spellChecker: false,
        shortcuts: {
          togglePreview: null,
          toggleSideBySide: null,
          toggleFullScreen: null
        }
      }"
    />

    <div class="mt-2 flex gap-2">
      <Button
        v-if="!hideCancel"
        variant="secondary"
        class="ml-auto"
        size="sm"
        @click="emit('cancel')"
      >
        Cancel
      </Button>
      <Button
        :variant="focus ? 'default' : 'secondary'"
        size="sm"
        @click="emit('save')"
        class="transition-colors"
      >
        Save
        <ShortCut keys="⌘ ↵" />
      </Button>
    </div>
  </div>
</template>

<style>
.EasyMDEContainer .CodeMirror {
  border: none;
  padding: 0;
  color: hsl(var(--card-foreground));
  background-color: hsl(var(--card));
}

.CodeMirror-cursor {
  border-color: hsl(var(--border));
}

@media (prefers-color-scheme: dark) {
  .CodeMirror-cursor {
    border-color: hsl(var(--foreground));
  }
}
</style>
