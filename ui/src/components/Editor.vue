<template>
  <prism-editor
      v-if="showEditor"
      class="my-editor"
      v-model="code"
      :highlight="highlighter"
      line-numbers
      :readonly="readonly">
  </prism-editor>

  <!--MonacoEditor
      v-if="showEditor"
      ref="editor"
      class="editor"
      style="height: 100%"
      v-model="code"
      :language="this.lang"
      :options="{ scrollBeyondLastLine: false }"
      :theme="$vuetify.theme.dark ? 'vs-dark' : 'vs'"
  /-->
</template>

<script lang="ts">
import Vue from "vue";
// import MonacoEditor from "vue-monaco";

import { PrismEditor } from 'vue-prism-editor';
import 'vue-prism-editor/dist/prismeditor.min.css'; // import the styles somewhere

import { highlight, languages } from 'prismjs/components/prism-core';
// import 'prismjs/components/prism-javascript';
import 'prismjs/components/prism-json';
import 'prismjs/components/prism-python';
import 'prismjs/components/prism-yaml';
import 'prismjs/components/prism-markup';
import 'prismjs/components/prism-log';
import 'prismjs/themes/prism-tomorrow.css'; // import syntax highlighting styles

interface State {
  showEditor: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  resize?: any;
  code: string,
}

export default Vue.extend({
  name: "Editor",
  props: ["value", "lang", "readonly"],
  components: {PrismEditor},
  data: (): State => ({
    showEditor: true,
    resize: undefined,
    code: "",
  }),
  watch: {
    code: function () {
      this.$emit('input', this.code);
    },
    value: function () {
      this.code = this.value;
    }
  },
  methods: {
    // resizeEditor() {
    //   this.showEditor = false;
    //   this.$nextTick(() => {
    //     this.showEditor = true;
    //   });
    // },
    highlighter(code: string) {
      switch (this.lang) {
        case "python":
          return highlight(code, languages.python);
        case "log":
          return highlight(code, languages.log);
        case "yaml":
          return highlight(code, languages.yaml);
        case "json":
          return highlight(code, languages.json);
        case "html":
          return highlight(code, languages.html);
      }
      return highlight(code, languages.json);
    },
  },
  mounted() {
    this.code = this.value;
    // this.resize = this.lodash.debounce(this.resizeEditor, 200);
    // window.addTicketListener("resize", this.resize);
  },
  destroyed() {
    // window.removeticketListener("resize", this.resize);
  },
});
</script>

<style>
/* required class */
.my-editor {
  /* we dont use `language-` classes anymore so thats why we need to add background and text color manually */
  background: #2d2d2d;
  color: #ccc;

  /* you must provide font-family font-size line-height. Example: */
  font-family: Fira code, Fira Mono, Consolas, Menlo, Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  padding: 5px;
}

/* optional class for removing the outline */
.prism-editor__textarea:focus {
  outline: none;
}
</style>
