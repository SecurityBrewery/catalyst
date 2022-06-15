<template>
  <div v-if="template !== undefined" class="flex-grow-1 flex-column d-flex fill-height pa-8">
    <v-row class="flex-grow-0 flex-shrink-0">
      <v-spacer></v-spacer>
      <v-btn href="https://catalyst-soar.com/docs/catalyst/engineer/template" target="_blank" outlined rounded small>
        <v-icon>mdi-book-open</v-icon> Handbook
      </v-btn>
    </v-row>

    <v-alert v-if="readonly" type="info">
      You do not have write access to templates.
      Changes here cannot be saved.
    </v-alert>
    <div class="d-flex" style="align-items: center">
      <h2 v-if="readonly">Template: {{ template.name }}</h2>
      <h2 v-else-if="this.$route.params.id === 'new'">New Template: {{ template.name }}</h2>
      <h2 v-else>Edit Template: {{ template.name }}</h2>
    </div>

    <v-text-field id="name-edit" label="Name" v-model="template.name" class="flex-grow-0 flex-shrink-0" :readonly="readonly"></v-text-field>

    <AdvancedJSONSchemaEditor id="template-edit" v-if="schema" @save="save" :schema="schema" :readonly="readonly" :hidepreview="false"></AdvancedJSONSchemaEditor>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import { TicketTemplate } from "@/client";
import AdvancedJSONSchemaEditor from "../components/AdvancedJSONSchemaEditor.vue";

interface State {
  template?: TicketTemplate,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  schema?: any;
}

export default Vue.extend({
  name: "Template",
  components: { AdvancedJSONSchemaEditor },
  data: (): State => ({
    template: undefined,
    schema: undefined,
  }),
  watch: {
    $route: 'loadTemplate',
  },
  computed: {
    readonly: function (): boolean {
      return !this.hasRole("engineer:template:write");
    },
  },
  methods: {
    save(schema) {
      if (this.template === undefined) {
        return;
      }
      let template = this.template;
      template.schema = schema;

      if (this.$route.params.id == 'new') {
        this.$store.dispatch("addTemplate", template).then(() => {
          this.$router.push({name: "TemplateList" });
        });
      } else {
        this.$store.dispatch("updateTemplate", { id: this.$route.params.id, template: template });
      }
    },
    loadTemplate() {
      if (!this.$route.params.id) {
        return
      }
      if (this.$route.params.id == 'new') {
        this.template = { name: "MyTemplate", schema: '{ "type": "object", "name": "Incident" }' }
        this.schema = { type: "object", name: "Incident" }
      } else {
        this.$store.dispatch("getTemplate", this.$route.params.id).then((response: TicketTemplate) => {
          this.template = response;
          this.schema = JSON.parse(response.schema);
        });
      }
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    }
  },
  mounted() {
    this.loadTemplate();
  },
});
</script>
